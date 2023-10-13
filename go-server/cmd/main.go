package main

import (
	"os"

	"github.com/acornak/poc-gpt/handlers"
	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/acornak/poc-gpt/docs"
)

type Logger interface {
	Info(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
}

type Server struct {
	Router  *gin.Engine
	Logger  Logger
	Handler *handlers.Handler
}

func NewServer(logger *zap.Logger) *Server {
	router := gin.Default()
	handler := &handlers.Handler{Logger: logger}
	s := &Server{Router: router, Logger: logger, Handler: handler}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.POST("/add", handler.Add)
	router.POST("/subtract", handler.Subtract)
	router.POST("/compute", handler.Compute)

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

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Failed to load env variables: %v\n", zap.Error(err))
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "debug"
		logger.Info("GIN_MODE not found in env, using 'debug' as default")
	}
	gin.SetMode(ginMode)

	server := NewServer(logger)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		logger.Info("PORT not found in env, using 8080 as default")
	}

	if err := server.Router.Run(":" + port); err != nil {
		server.Logger.Fatal("Couldn't start server: %v\n", zap.Error(err))
	}
}
