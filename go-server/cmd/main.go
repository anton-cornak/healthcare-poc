package main

import (
	"github.com/acornak/poc-gpt/handlers"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/acornak/poc-gpt/docs"
)

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

	handler := &handlers.Handler{Logger: logger}

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/add", handler.Add)
	r.POST("/subtract", handler.Subtract)
	r.POST("/compute", handler.Compute)

	err = r.Run(":8080")
	if err != nil {
		logger.Fatal("Couldn't start server: %v\n", zap.Error(err))
	}
}
