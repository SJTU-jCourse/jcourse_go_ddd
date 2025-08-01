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
	"jcourse_go/internal/domain/event"
	"jcourse_go/internal/interface/dto"
	"jcourse_go/internal/interface/task"
	"jcourse_go/internal/interface/web"
)

const (
	SuccessCode           = 0
	SuccessMessage        = "success"
	DefaultPort           = "8080"
	ShutdownTimeout       = 30 * time.Second
	WorkerCleanupDelay    = 5 * time.Second
	SignalChannelCapacity = 1
	HealthCheckEndpoint   = "/health"
	HTTPStatusOK          = 200
)

func loadConfiguration() *config.Config {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "" {
		gin.SetMode(gin.ReleaseMode)
	}

	return cfg
}

func setupServiceContainer(cfg *config.Config, eventPublisher event.Publisher) *app.ServiceContainer {
	serviceContainer, err := app.NewServiceContainer(*cfg, eventPublisher)
	if err != nil {
		log.Fatalf("Failed to initialize service container: %v", err)
	}
	return serviceContainer
}

func setupEventbus(cfg *config.Config, serviceContainer *app.ServiceContainer) *app.EventBusSetup {
	eventBusSetup, err := app.SetupEventBus(*cfg, serviceContainer)
	if err != nil {
		log.Fatalf("Failed to setup eventbus: %v", err)
	}
	return eventBusSetup
}

func startBackgroundWorkers(ctx context.Context, cfg *config.Config, eventBusSetup *app.EventBusSetup, serviceContainer *app.ServiceContainer) {
	if cfg.Event.Enabled {
		log.Println("Starting background event handlers...")

		// Start event bus worker (async event processing)
		go func() {
			if err := eventBusSetup.StartEventBus(); err != nil {
				log.Printf("Failed to start event bus: %v", err)
			}
		}()

		// Start email worker
		emailWorker := task.NewEmailWorker(serviceContainer)
		go emailWorker.Start(ctx)

		// Start statistics worker
		statsWorker := task.NewStatisticsWorker(serviceContainer)
		go statsWorker.Start(ctx)

		// Start cleanup worker
		cleanupWorker := task.NewCleanupWorker(serviceContainer)
		go cleanupWorker.Start(ctx)

		log.Println("Background workers started successfully")
	}
}

func setupHTTPServer(serviceContainer *app.ServiceContainer) *http.Server {
	// Initialize Gin router
	router := gin.New()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(web.CORSMiddleware())

	// Register routes
	web.RegisterRouter(router, serviceContainer)

	// Add health check endpoint
	router.GET(HealthCheckEndpoint, func(c *gin.Context) {
		response := dto.BaseResponse{
			Code: SuccessCode,
			Msg:  SuccessMessage,
		}
		c.JSON(HTTPStatusOK, response)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = DefaultPort
	}

	// Start server
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting unified server on %s", addr)
	log.Printf("Health check available at http://localhost%s/health", addr)

	// Create HTTP server
	return &http.Server{
		Addr:    addr,
		Handler: router,
	}
}

func waitForShutdown(server *http.Server, sigChan chan os.Signal, cancel context.CancelFunc) {
	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, gracefully stopping server...")

	// Cancel context to stop all background workers
	cancel()

	// Create shutdown context with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer shutdownCancel()

	// Gracefully shutdown the server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server stopped gracefully")
	}

	// Give workers time to clean up
	time.Sleep(WorkerCleanupDelay)
	log.Println("All services stopped")
}

func main() {
	// Load configuration
	cfg := loadConfiguration()

	// Initialize service container first
	serviceContainer := setupServiceContainer(cfg, nil)
	defer serviceContainer.Close()

	// Setup eventbus using service container
	eventBusSetup := setupEventbus(cfg, serviceContainer)
	defer eventBusSetup.ShutdownEventBus()

	// Create context with cancellation for background tasks
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, SignalChannelCapacity)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start background workers if event system is enabled
	startBackgroundWorkers(ctx, cfg, eventBusSetup, serviceContainer)

	// Setup HTTP server
	server := setupHTTPServer(serviceContainer)

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown
	waitForShutdown(server, sigChan, cancel)
}
