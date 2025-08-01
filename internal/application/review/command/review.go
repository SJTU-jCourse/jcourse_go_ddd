package command

import (
	"fmt"

	"jcourse_go/internal/domain/common"
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
	reviewRepo            review.ReviewRepository
	courseRepo            review.CourseRepository
	permissionService     permission.ReviewPermissionService
}

func NewReviewCommandService(
	reviewRepo review.ReviewRepository,
	courseRepo review.CourseRepository,
	permissionService permission.ReviewPermissionService) ReviewCommandService {
	return &reviewCommandService{
		reviewRepo:        reviewRepo,
		courseRepo:        courseRepo,
		permissionService: permissionService,
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
	// 2. todo 频控
	// 3. todo 内容校验
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
	action := review.NewReviewAction(reviewID, commonCtx.User.UserID, actionType)
	if err := s.reviewRepo.SaveReviewAction(commonCtx.Ctx, &action); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "post_review_action").WithMetadata("review_id", reviewID).WithMetadata("action_type", actionType)
	}
	return nil
}

func (s *reviewCommandService) DeleteReviewAction(commonCtx *common.CommonContext, reviewID int, actionID int) error {
	if err := s.reviewRepo.DeleteReviewAction(commonCtx.Ctx, actionID); err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "delete_review_action").WithMetadata("review_id", reviewID).WithMetadata("action_id", actionID)
	}
	return nil
}
