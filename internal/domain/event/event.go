package event

import "context"

type Event interface {
	ID() string
	Type() Type
	Payload() Payload
}

type Payload interface {
	Type() Type
}

type Type int

const (
	TypeUserCreated   Type = iota
	TypeReviewCreated Type = iota
)

type Publisher interface {
	Publish(ctx context.Context, events ...Event) error
}
