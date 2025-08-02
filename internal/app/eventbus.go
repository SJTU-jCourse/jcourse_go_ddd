package app

import (
	"strconv"

	"jcourse_go/internal/config"
	"jcourse_go/internal/domain/event"
	eventbusimpl "jcourse_go/internal/infrastructure/eventbus"
	eventhandler "jcourse_go/internal/interface/handler"
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

	redisAddr := conf.Redis.Addr
	if conf.Redis.Port != 0 {
		redisAddr = redisAddr + ":" + strconv.Itoa(conf.Redis.Port)
	}

	eventBus, err := eventbusimpl.NewAsynqEventBus(redisAddr)
	if err != nil {
		return nil, err
	}

	// Register event handlers
	if err := eventhandler.RegisterEventHandlers(eventBus, serviceContainer.GetPointCommandService()); err != nil {
		return nil, err
	}

	return &EventBusSetup{
		EventBus: eventBus,
		Enabled:  true,
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
