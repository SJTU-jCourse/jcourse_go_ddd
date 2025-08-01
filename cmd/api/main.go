package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/web"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize service container
	serviceContainer, err := app.NewServiceContainer(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize service container: %v", err)
	}
	defer serviceContainer.Close()

	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(web.CORSMiddleware())

	// Register routes
	web.RegisterRouter(router, serviceContainer)

	// Add health check endpoint
	router.GET("/health", func(c *gin.Context) {
		response := dto.BaseResponse{
			Code: 0,
			Msg:  "success",
		}
		c.JSON(200, response)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	log.Printf("Health check available at http://localhost%s/health", addr)

	// Create HTTP server
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, gracefully stopping server...")

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}
}
