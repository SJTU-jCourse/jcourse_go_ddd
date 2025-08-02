package event

import "context"

type Handler interface {
	Handle(ctx context.Context, e Event) error
}

type EventBus interface {
	Register(eventType Type, handler Handler) error
	Dispatch(ctx context.Context, events ...Event) error
	Start() error
	Shutdown() error
}

type EventBusPublisher interface {
	EventBus
	Publisher
}
