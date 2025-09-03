package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bankmore/internal/fee/domain"
	"bankmore/internal/fee/handlers"
	"bankmore/internal/fee/repository"
	"bankmore/internal/fee/service"
	"bankmore/internal/shared/kafka"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title BankMore Fee API
// @version 1.0
// @description API de tarifas do sistema bancário BankMore
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8003
// @BasePath /

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

	if err := db.AutoMigrate(&domain.Fee{}); err != nil {
		logger.WithError(err).Fatal("Failed to migrate database")
	}

	feeRepo := repository.NewFeeRepository(db)
	feeService := service.NewFeeService(feeRepo, logger)
	feeHandler := handlers.NewFeeHandler(feeService, logger)

	consumer, err := kafka.NewConsumer("fee-service", feeService, logger)
	if err != nil {
		logger.WithError(err).Fatal("Failed to create Kafka consumer")
	}
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		logger.Info("Starting Kafka consumer")
		if err := consumer.Start(ctx); err != nil {
			logger.WithError(err).Error("Kafka consumer error")
		}
	}()

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

	api := router.Group("/api/fee")
	{
		api.GET("/:accountNumber", feeHandler.GetFeesByAccount)
		api.GET("/fee/:id", feeHandler.GetFeeByID)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"service":   "fee-api",
			"timestamp": time.Now().UTC(),
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8003"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		logger.WithField("port", port).Info("Starting Fee API server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.WithError(err).Fatal("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down Fee API server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.WithError(err).Fatal("Server forced to shutdown")
	}

	logger.Info("Fee API server exited")
}
