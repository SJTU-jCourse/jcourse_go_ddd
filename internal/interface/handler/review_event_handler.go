package handler

import (
	"context"
	"fmt"
	"log"

	"jcourse_go/internal/domain/event"
)

type ReviewEventHandler struct{}

func NewReviewEventHandler() *ReviewEventHandler {
	return &ReviewEventHandler{}
}

func (h *ReviewEventHandler) Handle(ctx context.Context, e event.Event) error {
	payload, ok := e.Payload().(*event.ReviewPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for review event")
	}

	switch e.Type() {
	case event.TypeReviewCreated:
		log.Printf("Review created event: ReviewID=%d, UserID=%d, CourseID=%d, Rating=%d",
			payload.ReviewID, payload.UserID, payload.CourseID, payload.Rating)
	case event.TypeReviewModified:
		log.Printf("Review modified event: ReviewID=%d, UserID=%d, CourseID=%d, Rating=%d",
			payload.ReviewID, payload.UserID, payload.CourseID, payload.Rating)
	default:
		return fmt.Errorf("unsupported review event type: %d", e.Type())
	}

	return nil
}
