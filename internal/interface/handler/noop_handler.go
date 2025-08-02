package handler

import (
	"context"

	"jcourse_go/internal/domain/event"
)

type NoOpHandler struct{}

func NewNoOpHandler() *NoOpHandler {
	return &NoOpHandler{}
}

func (h *NoOpHandler) Handle(ctx context.Context, e event.Event) error {
	return nil
}