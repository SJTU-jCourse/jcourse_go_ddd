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
		return apperror.ErrDB.Wrap(err)
	}
	if c == nil {
		return apperror.ErrNoTargetCourse
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
		return err
	}
	return nil
}

func (s *reviewCommandService) UpdateReview(commonCtx *common.CommonContext, cmd *review.UpdateReviewCommand) error {
	r, err := s.reviewRepo.Get(commonCtx.Ctx, cmd.ReviewID)
	if err != nil {
		return err
	}
	if r == nil {
		return nil
	}

	// Check permission
	canUpdate, reason := s.permissionService.CanUpdateReview(commonCtx.Ctx, permission.ToPermissionReview(r), commonCtx.User)
	if !canUpdate {
		return apperror.ErrPermission.WithMessage(fmt.Sprintf("cannot update review: %s", reason))
	}

	revision := review.NewRevisionFromReview(r)
	r.Update(&cmd.ReviewContent)
	if err := s.ValidateReview(commonCtx, r); err != nil {
		return err
	}
	if err := s.reviewRepo.Save(commonCtx.Ctx, r, &revision); err != nil {
		return err
	}
	return nil
}

func (s *reviewCommandService) DeleteReview(commonCtx *common.CommonContext, cmd *review.DeleteReviewCommand) error {
	r, err := s.reviewRepo.Get(commonCtx.Ctx, cmd.ReviewID)
	if err != nil {
		return err
	}
	if r == nil {
		return nil
	}

	// Check permission
	canDelete, reason := s.permissionService.CanDeleteReview(commonCtx.Ctx, permission.ToPermissionReview(r), commonCtx.User)
	if !canDelete {
		return apperror.ErrPermission.WithMessage(fmt.Sprintf("cannot delete review: %s", reason))
	}

	if err := s.reviewRepo.Delete(commonCtx.Ctx, review.ReviewFilter{ReviewID: &cmd.ReviewID}); err != nil {
		return err
	}
	return nil
}

func (s *reviewCommandService) PostReviewAction(commonCtx *common.CommonContext, reviewID int, actionType string) error {
	action := review.NewReviewAction(reviewID, commonCtx.User.UserID, actionType)
	return s.reviewRepo.SaveReviewAction(commonCtx.Ctx, &action)
}

func (s *reviewCommandService) DeleteReviewAction(commonCtx *common.CommonContext, reviewID int, actionID int) error {
	return s.reviewRepo.DeleteReviewAction(commonCtx.Ctx, actionID)
}
