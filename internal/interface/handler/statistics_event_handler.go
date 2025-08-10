package handler

import (
	"context"
	"fmt"
	"log"

	"jcourse_go/internal/domain/event"
)

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
		log.Printf("Updating statistics for new review: CourseID=%d", payload.CourseID)
	case event.TypeReviewModified:
		log.Printf("Updating statistics for modified review: CourseID=%d", payload.CourseID)
	}

	return nil
}
