package main

import (
	"context"
	"log"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/interface/task"
)

func StartBackgroundWorkers(ctx context.Context, cfg *config.Config, eventBusSetup *app.EventBusSetup, serviceContainer *app.ServiceContainer) {
	if cfg.Event.Enabled {
		log.Println("Starting background event handlers...")

		// Start event bus worker (async event processing)
		go func() {
			if err := eventBusSetup.StartEventBus(); err != nil {
				log.Printf("Failed to start event bus: %v", err)
			}
		}()

		// Start statistics worker
		statsWorker := task.NewStatisticsWorker(serviceContainer)
		go statsWorker.Start(ctx)

		log.Println("Background workers started successfully")
	}
}