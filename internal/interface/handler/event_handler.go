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
		log.Printf("Review created event: ReviewID=%s, UserID=%s, CourseID=%s, Rating=%d",
			payload.ReviewID, payload.UserID, payload.CourseID, payload.Rating)
	case event.TypeReviewModified:
		log.Printf("Review modified event: ReviewID=%s, UserID=%s, CourseID=%s, Rating=%d",
			payload.ReviewID, payload.UserID, payload.CourseID, payload.Rating)
	default:
		return fmt.Errorf("unsupported review event type: %d", e.Type())
	}

	return nil
}

type PointEventHandler struct{}

func NewPointEventHandler() *PointEventHandler {
	return &PointEventHandler{}
}

func (h *PointEventHandler) Handle(ctx context.Context, e event.Event) error {
	payload, ok := e.Payload().(*event.ReviewPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for point event")
	}

	if e.Type() == event.TypeReviewCreated {
		log.Printf("Awarding points for review creation: UserID=%s, ReviewID=%s",
			payload.UserID, payload.ReviewID)
	}

	return nil
}

type StatisticsEventHandler struct{}

func NewStatisticsEventHandler() *StatisticsEventHandler {
	return &StatisticsEventHandler{}
}

func (h *StatisticsEventHandler) Handle(ctx context.Context, e event.Event) error {
	payload, ok := e.Payload().(*event.ReviewPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for statistics event")
	}

	switch e.Type() {
	case event.TypeReviewCreated:
		log.Printf("Updating statistics for new review: CourseID=%s", payload.CourseID)
	case event.TypeReviewModified:
		log.Printf("Updating statistics for modified review: CourseID=%s", payload.CourseID)
	}

	return nil
}

type NoOpHandler struct{}

func NewNoOpHandler() *NoOpHandler {
	return &NoOpHandler{}
}

func (h *NoOpHandler) Handle(ctx context.Context, e event.Event) error {
	return nil
}
