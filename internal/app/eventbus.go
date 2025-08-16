package app

import (
	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/event"
)

// EventBusSetup holds the eventbus and its configuration
type EventBusSetup struct {
	EventBus event.EventBusPublisher
	Enabled  bool
}

// SetupEventBus creates and configures the eventbus, but doesn't start it
func SetupEventBus(conf config.Config, serviceContainer *ServiceContainer) (*EventBusSetup, error) {
	if !conf.Event.Enabled {
		return &EventBusSetup{
			EventBus: nil,
			Enabled:  false,
		}, nil
	}

	// For now, return a disabled eventbus since we removed asynq
	// The domain interfaces are preserved for future implementation
	return &EventBusSetup{
		EventBus: nil,
		Enabled:  false,
	}, nil
}

// StartEventBus starts the eventbus worker
func (e *EventBusSetup) StartEventBus() error {
	if e.EventBus != nil {
		return e.EventBus.Start()
	}
	return nil
}

// ShutdownEventBus shuts down the eventbus
func (e *EventBusSetup) ShutdownEventBus() {
	if e.EventBus != nil {
		e.EventBus.Shutdown()
	}
}

// GetPublisher returns the event publisher for services
func (e *EventBusSetup) GetPublisher() event.Publisher {
	if e.EventBus != nil {
		return e.EventBus
	}
	return nil
}
