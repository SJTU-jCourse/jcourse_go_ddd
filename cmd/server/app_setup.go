package main

import (
	"log"

	"jcourse_go/internal/app"
	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/event"
)

func SetupServiceContainer(cfg *config.Config, eventPublisher event.Publisher) *app.ServiceContainer {
	serviceContainer, err := app.NewServiceContainer(*cfg, eventPublisher)
	if err != nil {
		log.Fatalf("Failed to initialize service container: %v", err)
	}
	return serviceContainer
}

func SetupEventbus(cfg *config.Config, serviceContainer *app.ServiceContainer) *app.EventBusSetup {
	eventBusSetup, err := app.SetupEventBus(*cfg, serviceContainer)
	if err != nil {
		log.Fatalf("Failed to setup eventbus: %v", err)
	}
	return eventBusSetup
}