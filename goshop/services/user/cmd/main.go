package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/goshop/pkg/config"
	"github.com/yourusername/goshop/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const serviceName = "user"

func main() {
	// Load configuration
	cfg, err := config.Load(serviceName, "")
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log, err := logger.New(serviceName, cfg.Service.LogLevel)
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer log.Sync()

	ctx := context.Background()
	log.Info(ctx, "Starting user service",
		zap.String("environment", cfg.Service.Environment),
		zap.Int("http_port", cfg.HTTP.Port),
		zap.Int("grpc_port", cfg.GRPC.Port),
	)

	// Initialize HTTP server
	router := gin.Default()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTP.Port),
		Handler: router,
	}

	// Register HTTP routes
	setupHTTPRoutes(router)

	// Initialize gRPC server
	grpcServer := grpc.NewServer()
	// Register gRPC services

	// Start HTTP server
	go func() {
		log.Info(ctx, "Starting HTTP server", zap.Int("port", cfg.HTTP.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(ctx, "HTTP server failed", zap.Error(err))
		}
	}()

	// Start gRPC server
	go func() {
		log.Info(ctx, "Starting gRPC server", zap.Int("port", cfg.GRPC.Port))
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPC.Port))
		if err != nil {
			log.Fatal(ctx, "Failed to listen on gRPC port", zap.Error(err))
		}
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(ctx, "gRPC server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info(ctx, "Received shutdown signal")

	// Gracefully shutdown servers
	log.Info(ctx, "Shutting down servers")
	grpcServer.GracefulStop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatal(ctx, "Server forced to shutdown", zap.Error(err))
	}

	log.Info(ctx, "Server has been shutdown successfully")
}

// Setup HTTP routes
func setupHTTPRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "UP",
		})
	})

	api := router.Group("/api/v1")
	{
		users := api.Group("/users")
		{
			users.POST("/register", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.POST("/login", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.POST("/reset-password", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.GET("/me", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.PUT("/me", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.GET("/me/addresses", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})

			users.POST("/me/addresses", func(c *gin.Context) {
				// Not implemented yet
				c.JSON(http.StatusOK, gin.H{"message": "Not implemented"})
			})
		}
	}
}
