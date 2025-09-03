package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bankmore/internal/shared/kafka"
	"bankmore/internal/shared/middleware"
	"bankmore/internal/transfer/domain"
	"bankmore/internal/transfer/handlers"
	"bankmore/internal/transfer/repository"
	"bankmore/internal/transfer/service"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title BankMore Transfer API
// @version 1.0
// @description API de transferências do sistema bancário BankMore
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8002
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./database/bankmore.db"
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.WithError(err).Fatal("Failed to connect to database")
	}

	if err := db.AutoMigrate(&domain.Transfer{}); err != nil {
		logger.WithError(err).Fatal("Failed to migrate database")
	}

	producer, err := kafka.NewProducer(logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Kafka producer")
	}
	defer producer.Close()

	transferRepo := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(transferRepo, producer, logger)
	transferHandler := handlers.NewTransferHandler(transferService, logger)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	api := router.Group("/api/transfer")
	api.Use(middleware.JWTMiddleware())
	{
		api.POST("", transferHandler.CreateTransfer)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "transfer-api",
			"timestamp": time.Now().UTC(),
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.WithField("port", port).Info("Starting Transfer API server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Transfer API server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Transfer API server exited")
}
