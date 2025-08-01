package permission

import (
	"context"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
)

type ReviewPermissionService interface {
	CanUpdateReview(ctx context.Context, review *Review, user *common.User) (bool, string)
	CanDeleteReview(ctx context.Context, review *Review, user *common.User) (bool, string)
}

type ReviewPermissionChecker struct {
	userRepo auth.UserRepository
}

func NewReviewPermissionChecker(userRepo auth.UserRepository) ReviewPermissionService {
	return &ReviewPermissionChecker{
		userRepo: userRepo,
	}
}

func (p *ReviewPermissionChecker) CanUpdateReview(ctx context.Context, review *Review, user *common.User) (bool, string) {
	return p.checkReviewPermission(ctx, review, user, ActionUpdate)
}

func (p *ReviewPermissionChecker) CanDeleteReview(ctx context.Context, review *Review, user *common.User) (bool, string) {
	return p.checkReviewPermission(ctx, review, user, ActionDelete)
}

func (p *ReviewPermissionChecker) checkReviewPermission(ctx context.Context, review *Review, user *common.User, action Action) (bool, string) {
	if user == nil {
		return false, "not authenticated"
	}

	// Check if user is admin
	if user.Role == common.RoleAdmin {
		return true, "admin access"
	}

	// Check if user is the owner
	if review.UserID == user.UserID {
		return true, "owner access"
	}

	return false, "permission denied"
}

// Review is a lightweight struct for permission checking
type Review struct {
	ID     int
	UserID int
}

// Helper function to convert domain review to permission review
func ToPermissionReview(review *review.Review) *Review {
	return &Review{
		ID:     review.ID,
		UserID: review.UserID,
	}
}
