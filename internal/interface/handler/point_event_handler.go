package handler

import (
	"context"
	"fmt"
	"log"

	"jcourse_go/internal/application/point/command"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/event"
)

type PointEventHandler struct {
	pointService command.PointCommandService
}

func NewPointEventHandler(pointService command.PointCommandService) *PointEventHandler {
	return &PointEventHandler{
		pointService: pointService,
	}
}

func (h *PointEventHandler) Handle(ctx context.Context, e event.Event) error {
	payload, ok := e.Payload().(*event.ReviewPayload)
	if !ok {
		return fmt.Errorf("invalid payload type for point event")
	}

	if e.Type() == event.TypeReviewCreated {
		log.Printf("Awarding points for review creation: UserID=%d, ReviewID=%d",
			payload.UserID, payload.ReviewID)

		// Award points using the service layer
		return h.awardPointsForReview(ctx, payload.UserID, payload.ReviewID)
	}

	return nil
}

// awardPointsForReview awards points to a user for creating a review using the service layer
func (h *PointEventHandler) awardPointsForReview(ctx context.Context, userID int, reviewID int) error {
	// Create system context for internal operations
	systemCtx := &common.CommonContext{
		Ctx:  ctx,
		User: common.SystemUser, // Use system user for internal operations
	}

	// Use the point service to award points
	if err := h.pointService.AwardPointsForReview(systemCtx, userID, reviewID); err != nil {
		log.Printf("Failed to award points for review %d to user %d: %v", reviewID, userID, err)
		return err
	}

	log.Printf("Successfully awarded points to user %d for review %d", userID, reviewID)
	return nil
}
