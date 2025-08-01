package event

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event interface {
	ID() string
	Type() Type
	Payload() Payload
	ToJSON() ([]byte, error)
}

type Payload interface {
	Type() Type
}

type Type int

const (
	TypeUserCreated    Type = iota
	TypeReviewCreated  Type = iota
	TypeReviewModified Type = iota
)

type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
}

type BaseEvent struct {
	id        string
	eventType Type
	payload   Payload
	timestamp time.Time
}

func NewBaseEvent(eventType Type, payload Payload) *BaseEvent {
	return &BaseEvent{
		id:        uuid.New().String(),
		eventType: eventType,
		payload:   payload,
		timestamp: time.Now(),
	}
}

func (e *BaseEvent) ID() string {
	return e.id
}

func (e *BaseEvent) Type() Type {
	return e.eventType
}

func (e *BaseEvent) Payload() Payload {
	return e.payload
}

func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

func (e *BaseEvent) ToJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":        e.id,
		"type":      e.eventType,
		"payload":   e.payload,
		"timestamp": e.timestamp,
	})
}

func (e *BaseEvent) SetPayload(payload Payload) {
	e.payload = payload
}

func (e *BaseEvent) EventType() Type {
	return e.eventType
}
