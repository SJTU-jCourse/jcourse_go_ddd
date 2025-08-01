package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/interface/task"
)

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize service container
	serviceContainer, err := app.NewServiceContainer(*cfg)
	if err != nil {
		log.Fatalf("Failed to initialize service container: %v", err)
	}
	defer serviceContainer.Close()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start workers
	log.Println("Starting background workers...")

	// Start email worker
	emailWorker := task.NewEmailWorker(serviceContainer)
	go emailWorker.Start(ctx)

	// Start statistics worker
	statsWorker := task.NewStatisticsWorker(serviceContainer)
	go statsWorker.Start(ctx)

	// Start cleanup worker
	cleanupWorker := task.NewCleanupWorker(serviceContainer)
	go cleanupWorker.Start(ctx)

	log.Println("Worker service started successfully")

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, gracefully stopping workers...")

	// Cancel context to stop all workers
	cancel()

	// Give workers time to clean up
	time.Sleep(5 * time.Second)

	log.Println("Worker service stopped")
}
