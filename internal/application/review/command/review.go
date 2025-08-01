package command

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/event"
	"jcourse_go/internal/domain/permission"
	"jcourse_go/internal/domain/review"
	"jcourse_go/pkg/apperror"
)

type ReviewCommandService interface {
	WriteReview(commonCtx *common.CommonContext, cmd *review.WriteReviewCommand) error
	UpdateReview(commonCtx *common.CommonContext, cmd *review.UpdateReviewCommand) error
	DeleteReview(commonCtx *common.CommonContext, cmd *review.DeleteReviewCommand) error
	PostReviewAction(commonCtx *common.CommonContext, reviewID int, actionType string) error
	DeleteReviewAction(commonCtx *common.CommonContext, reviewID int, actionID int) error
}

type reviewCommandService struct {
	reviewRepo        review.ReviewRepository
	courseRepo        review.CourseRepository
	permissionService permission.ReviewPermissionService
	eventPublisher    event.Publisher
}

func NewReviewCommandService(
	reviewRepo review.ReviewRepository,
	courseRepo review.CourseRepository,
	permissionService permission.ReviewPermissionService,
	eventPublisher event.Publisher) ReviewCommandService {
	return &reviewCommandService{
		reviewRepo:        reviewRepo,
		courseRepo:        courseRepo,
		permissionService: permissionService,
		eventPublisher:    eventPublisher,
	}
}

func (s *reviewCommandService) ValidateReview(commonCtx *common.CommonContext, r *review.Review) error {
	// 1. 课程 id 有效
	c, err := s.courseRepo.FindOfferedCourse(commonCtx.Ctx, r.CourseID, r.Semester)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "validate_review").WithMetadata("course_id", r.CourseID)
	}
	if c == nil {
		return apperror.ErrNoTargetCourse.WithMetadata("course_id", r.CourseID).WithMetadata("semester", r.Semester.String())
	}
	// 2. 频控：1分钟内最多发3条
	if err := s.checkRateLimit(commonCtx, r.UserID); err != nil {
		return err
	}
	// 3. 内容校验：最近3条内容相似度不能超过90%
	if err := s.checkContentSimilarity(commonCtx, r.UserID, r.Comment); err != nil {
		return err
	}
	return nil
}

func (s *reviewCommandService) WriteReview(commonCtx *common.CommonContext, cmd *review.WriteReviewCommand) error {
	r := review.NewReview(cmd.CourseID, commonCtx.User.UserID, &cmd.ReviewContent)
	if err := s.ValidateReview(commonCtx, &r); err != nil {
		return err
	}
	if err := s.reviewRepo.Save(commonCtx.Ctx, &r, nil); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "write_review").WithMetadata("user_id", commonCtx.User.UserID)
	}

	if s.eventPublisher != nil {
		payload := &event.ReviewPayload{
			ReviewID: strconv.Itoa(r.ID),
			UserID:   strconv.Itoa(r.UserID),
			CourseID: strconv.Itoa(r.CourseID),
			Rating:   r.Rating.Int(),
			Content:  r.Comment,
			Action:   "created",
		}

		reviewEvent := event.NewBaseEvent(event.TypeReviewCreated, payload)
		if err := s.eventPublisher.Publish(commonCtx.Ctx, reviewEvent); err != nil {
			return apperror.WrapDB(err).WithMetadata("operation", "publish_review_created_event").WithMetadata("review_id", r.ID)
		}
	}

	return nil
}

func (s *reviewCommandService) UpdateReview(commonCtx *common.CommonContext, cmd *review.UpdateReviewCommand) error {
	r, err := s.reviewRepo.Get(commonCtx.Ctx, cmd.ReviewID)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "update_review").WithMetadata("review_id", cmd.ReviewID)
	}
	if r == nil {
		return apperror.ErrNotFound.WithMessage("review not found").WithMetadata("review_id", cmd.ReviewID)
	}

	// Check permission
	canUpdate, reason := s.permissionService.CanUpdateReview(commonCtx.Ctx, permission.ToPermissionReview(r), commonCtx.User)
	if !canUpdate {
		return apperror.ErrPermission.WithMessage(fmt.Sprintf("cannot update review: %s", reason)).
			WithMetadata("review_id", cmd.ReviewID).
			WithMetadata("user_id", commonCtx.User.UserID).
			WithMetadata("owner_id", r.UserID)
	}

	revision := review.NewRevisionFromReview(r)
	r.Update(&cmd.ReviewContent)
	if err := s.ValidateReview(commonCtx, r); err != nil {
		return err
	}
	if err := s.reviewRepo.Save(commonCtx.Ctx, r, &revision); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "update_review").WithMetadata("review_id", cmd.ReviewID)
	}

	if s.eventPublisher != nil {
		payload := &event.ReviewPayload{
			ReviewID: strconv.Itoa(r.ID),
			UserID:   strconv.Itoa(r.UserID),
			CourseID: strconv.Itoa(r.CourseID),
			Rating:   r.Rating.Int(),
			Content:  r.Comment,
			Action:   "modified",
		}

		reviewEvent := event.NewBaseEvent(event.TypeReviewModified, payload)
		if err := s.eventPublisher.Publish(commonCtx.Ctx, reviewEvent); err != nil {
			return apperror.WrapDB(err).WithMetadata("operation", "publish_review_modified_event").WithMetadata("review_id", r.ID)
		}
	}

	return nil
}

func (s *reviewCommandService) DeleteReview(commonCtx *common.CommonContext, cmd *review.DeleteReviewCommand) error {
	r, err := s.reviewRepo.Get(commonCtx.Ctx, cmd.ReviewID)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "delete_review").WithMetadata("review_id", cmd.ReviewID)
	}
	if r == nil {
		return apperror.ErrNotFound.WithMessage("review not found").WithMetadata("review_id", cmd.ReviewID)
	}

	// Check permission
	canDelete, reason := s.permissionService.CanDeleteReview(commonCtx.Ctx, permission.ToPermissionReview(r), commonCtx.User)
	if !canDelete {
		return apperror.ErrPermission.WithMessage(fmt.Sprintf("cannot delete review: %s", reason)).
			WithMetadata("review_id", cmd.ReviewID).
			WithMetadata("user_id", commonCtx.User.UserID).
			WithMetadata("owner_id", r.UserID)
	}

	if err := s.reviewRepo.Delete(commonCtx.Ctx, review.ReviewFilter{ReviewID: &cmd.ReviewID}); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "delete_review").WithMetadata("review_id", cmd.ReviewID)
	}
	return nil
}

func (s *reviewCommandService) PostReviewAction(commonCtx *common.CommonContext, reviewID int, actionType string) error {
	// Check if user is authenticated
	if commonCtx.User.UserID == 0 {
		return apperror.ErrPermission.WithMessage("user not authenticated").
			WithMetadata("review_id", reviewID).
			WithMetadata("action_type", actionType)
	}

	// Check if review exists
	r, err := s.reviewRepo.Get(commonCtx.Ctx, reviewID)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "post_review_action").WithMetadata("review_id", reviewID)
	}
	if r == nil {
		return apperror.ErrNotFound.WithMessage("review not found").WithMetadata("review_id", reviewID)
	}

	action := review.NewReviewAction(reviewID, commonCtx.User.UserID, actionType)
	if err := s.reviewRepo.SaveReviewAction(commonCtx.Ctx, &action); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "post_review_action").WithMetadata("review_id", reviewID).WithMetadata("action_type", actionType)
	}
	return nil
}

func (s *reviewCommandService) DeleteReviewAction(commonCtx *common.CommonContext, reviewID int, actionID int) error {
	// Check if user is authenticated
	if commonCtx.User.UserID == 0 {
		return apperror.ErrPermission.WithMessage("user not authenticated").
			WithMetadata("review_id", reviewID).
			WithMetadata("action_id", actionID)
	}

	// Check if user owns the review action or is an admin
	action, err := s.reviewRepo.GetReviewAction(commonCtx.Ctx, actionID)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "get_review_action").WithMetadata("review_id", reviewID).WithMetadata("action_id", actionID)
	}
	if action == nil {
		return apperror.ErrNotFound.WithMessage("review action not found").
			WithMetadata("review_id", reviewID).
			WithMetadata("action_id", actionID)
	}

	if commonCtx.User.Role != common.RoleAdmin && action.UserID != commonCtx.User.UserID {
		return apperror.ErrPermission.WithMessage("user can only delete their own review actions").
			WithMetadata("review_id", reviewID).
			WithMetadata("action_id", actionID).
			WithMetadata("user_id", commonCtx.User.UserID).
			WithMetadata("action_owner_id", action.UserID)
	}

	if err := s.reviewRepo.DeleteReviewAction(commonCtx.Ctx, actionID); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "delete_review_action").WithMetadata("review_id", reviewID).WithMetadata("action_id", actionID)
	}
	return nil
}

func (s *reviewCommandService) checkRateLimit(commonCtx *common.CommonContext, userID int) error {
	// Find reviews created in the last minute
	oneMinuteAgo := time.Now().Add(-time.Minute)
	filter := review.ReviewFilter{
		UserID: &userID,
	}

	reviews, err := s.reviewRepo.FindBy(commonCtx.Ctx, filter)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "check_rate_limit").WithMetadata("user_id", userID)
	}

	// Count reviews created in the last minute
	recentCount := 0
	for _, r := range reviews {
		if r.CreatedAt.After(oneMinuteAgo) {
			recentCount++
		}
	}

	if recentCount >= 3 {
		return apperror.ErrRateLimit.WithMessage("rate limit exceeded: maximum 3 reviews per minute").
			WithMetadata("user_id", userID).
			WithMetadata("recent_count", recentCount).
			WithMetadata("limit", 3)
	}

	return nil
}

func (s *reviewCommandService) checkContentSimilarity(commonCtx *common.CommonContext, userID int, content string) error {
	// Find user's recent reviews
	filter := review.ReviewFilter{
		UserID: &userID,
	}

	reviews, err := s.reviewRepo.FindBy(commonCtx.Ctx, filter)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "check_content_similarity").WithMetadata("user_id", userID)
	}

	// Sort by creation time (most recent first)
	sort.Slice(reviews, func(i, j int) bool {
		return reviews[i].CreatedAt.After(reviews[j].CreatedAt)
	})

	// Check similarity with last 3 reviews
	maxCheck := min(3, len(reviews))
	for i := 0; i < maxCheck; i++ {
		similarity := calculateSimilarity(content, reviews[i].Comment)
		if similarity > 0.9 {
			return apperror.ErrValidation.WithMessage("content similarity too high with recent review").
				WithMetadata("user_id", userID).
				WithMetadata("similarity", similarity).
				WithMetadata("max_similarity", 0.9).
				WithMetadata("compared_review_id", reviews[i].ID)
		}
	}

	return nil
}

func calculateSimilarity(s1, s2 string) float64 {
	if s1 == "" || s2 == "" {
		return 0.0
	}

	// Simple similarity calculation based on common characters
	// This is a basic implementation - in production you might want to use more sophisticated algorithms
	distance := levenshteinDistance(s1, s2)
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 0.0
	}

	return 1.0 - float64(distance)/float64(maxLen)
}

func levenshteinDistance(s1, s2 string) int {
	r1, r2 := []rune(s1), []rune(s2)
	m, n := len(r1), len(r2)

	// Create distance matrix
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// Initialize first row and column
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}

	// Fill the matrix
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if r1[i-1] == r2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}
