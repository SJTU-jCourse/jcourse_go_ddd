package command

import (
	"fmt"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/point"
	"jcourse_go/pkg/apperror"
)

type PointCommandService interface {
	CreatePoint(commonCtx *common.CommonContext, userID int, amount int, reason string) error
	Transaction(commonCtx *common.CommonContext, fromUserID int, toUserID int, amount int, reason string) error
	AwardPointsForReview(commonCtx *common.CommonContext, userID int, reviewID int) error
}

func NewPointCommandService(repo point.UserPointRepository) PointCommandService {
	return &pointCommandService{
		repo: repo,
	}
}

type pointCommandService struct {
	repo point.UserPointRepository
}

func (s *pointCommandService) CreatePoint(commonCtx *common.CommonContext, userID int, amount int, reason string) error {
	// Check if user is admin
	if commonCtx.User.Role != common.RoleAdmin {
		return apperror.ErrPermission.WithMessage("only admins can create points").
			WithMetadata("user_id", userID).
			WithMetadata("amount", amount).
			WithMetadata("reason", reason)
	}

	// Create point record
	pointRecord := point.NewUserPointRecord(userID, amount, reason)
	if err := s.repo.Save(commonCtx.Ctx, &pointRecord); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "create_point").
			WithMetadata("user_id", userID).
			WithMetadata("amount", amount)
	}

	return nil
}

func (s *pointCommandService) Transaction(commonCtx *common.CommonContext, fromUserID int, toUserID int, amount int, reason string) error {
	// Check if user is admin
	if commonCtx.User.Role != common.RoleAdmin {
		return apperror.ErrPermission.WithMessage("only admins can perform point transactions").
			WithMetadata("from_user_id", fromUserID).
			WithMetadata("to_user_id", toUserID).
			WithMetadata("amount", amount).
			WithMetadata("reason", reason)
	}

	// Create transaction records
	fromPointRecord := point.NewUserPointRecord(fromUserID, -amount, reason)
	toPointRecord := point.NewUserPointRecord(toUserID, amount, reason)

	// Save both records
	if err := s.repo.Save(commonCtx.Ctx, &fromPointRecord); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "point_transaction_from").
			WithMetadata("from_user_id", fromUserID).
			WithMetadata("amount", amount)
	}

	if err := s.repo.Save(commonCtx.Ctx, &toPointRecord); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "point_transaction_to").
			WithMetadata("to_user_id", toUserID).
			WithMetadata("amount", amount)
	}

	return nil
}

// AwardPointsForReview awards points to a user for creating a review (internal system operation)
func (s *pointCommandService) AwardPointsForReview(commonCtx *common.CommonContext, userID int, reviewID int) error {
	// Define points to award for review creation
	const reviewCreationPoints = 10

	// Create point record with description
	reason := fmt.Sprintf("Points awarded for creating review #%d", reviewID)
	pointRecord := point.NewUserPointRecord(userID, reviewCreationPoints, reason)

	// Save the point record
	if err := s.repo.Save(commonCtx.Ctx, &pointRecord); err != nil {
		return apperror.WrapDB(err).
			WithMetadata("operation", "award_points_for_review").
			WithMetadata("user_id", userID).
			WithMetadata("review_id", reviewID).
			WithMetadata("points", reviewCreationPoints)
	}

	return nil
}
