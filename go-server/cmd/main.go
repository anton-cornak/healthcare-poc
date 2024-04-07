package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/acornak/healthcare-poc/handlers"
	"github.com/acornak/healthcare-poc/models"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	_ "github.com/acornak/healthcare-poc/docs"
	_ "github.com/lib/pq"
)

type logger interface {
	Info(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type server struct {
	Router  *gin.Engine
	Logger  logger
	Handler *handlers.Handler
}

type config struct {
	port   string
	env    string
	dbConn dbConfig
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
	sslmode  string
}

var apiVersion = "v1"

func loadConfigFromEnv(cfg *config) error {
	cfg.port = os.Getenv("PORT")
	cfg.dbConn.host = os.Getenv("DB_HOST")
	cfg.dbConn.port = os.Getenv("DB_PORT")
	cfg.dbConn.user = os.Getenv("DB_USER")
	cfg.dbConn.password = os.Getenv("DB_PASS")
	cfg.dbConn.dbname = os.Getenv("DB_NAME")
	cfg.dbConn.sslmode = os.Getenv("SSL_MODE")

	return validateConfig(cfg)
}

func initializeDatabase(cfg *dbConfig, logger *zap.Logger) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.host, cfg.port, cfg.user, cfg.password, cfg.dbname, cfg.sslmode)

	var db *sql.DB
	var err error

	deadline := time.Now().Add(5 * time.Minute)

	for attempts := 0; time.Now().Before(deadline); attempts++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			logger.Error("Error opening database connection", zap.Error(err))
			time.Sleep(time.Second * 5)
			continue
		}

		err = db.Ping()
		if err != nil {
			logger.Error("Database ping failed", zap.String("attempt", fmt.Sprintf("%d", attempts+1)), zap.Error(err))
			db.Close()
			time.Sleep(time.Second * 2)
		} else {
			logger.Info("Successfully connected to the database")
			return db, nil
		}
	}

	return nil, errors.New("failed to connect to the database")
}

func newServer(logger *zap.Logger, handler *handlers.Handler) *server {
	router := gin.Default()
	s := &server{Router: router, Logger: logger, Handler: handler}

	prefix := "/api/" + apiVersion

	// Documentation
	router.GET(prefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Math
	router.POST(prefix+"/math/add", handler.Add)
	router.POST(prefix+"/math/subtract", handler.Subtract)
	router.POST(prefix+"/math/compute", handler.Compute)

	// Location
	router.POST(prefix+"/location/wkt", handler.GetWKTLocation)
	router.POST(prefix+"/location/address", handler.GetAddressFromWKT)

	// Specialist
	router.POST(prefix+"/specialist/find", handler.FindSpecialist)
	// TODO
	// closest specialist
	// all specialists in area

	// Specialties
	router.POST(prefix+"/specialty/all", handler.GetSpecialties)

	return s
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		logger.Fatal("Couldn't initialize logger: %v\n", zap.Error(err))
	}

	defer func() {
		if err := logger.Sync(); err != nil {
			logger.Fatal("Error syncing logger: %v\n", zap.Error(err))
		}
	}()

	var cfg config
	flag.StringVar(&cfg.env, "ENV", "develop", "Application environment")
	flag.Parse()

	if cfg.env == "develop" {
		err := godotenv.Load(".envrc")
		if err != nil {
			fmt.Println("failed to load env vars:", err)
			os.Exit(1)
		}
	}

	err = loadConfigFromEnv(&cfg)
	if err != nil {
		fmt.Println("failed to load config:", err)
		os.Exit(1)
	}

	db, err := initializeDatabase(&cfg.dbConn, logger)
	if err != nil {
		logger.Fatal("failed to connect to DB:", zap.Error(err))
	}
	defer db.Close()

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
		logger.Info("GIN_MODE not found in env, using 'debug' as default")
	}
	gin.SetMode(ginMode)

	s := newServer(logger, &handlers.Handler{
		Logger: logger,
		Models: models.NewModels(db),
		Get:    http.Get,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("PORT not found in env, using 8080 as default")
	}

	if err := s.Router.Run(":" + port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
