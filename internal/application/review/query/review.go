package query

import (
	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
	"jcourse_go/pkg/apperror"
)

type ReviewQueryService interface {
	LatestReviews(commonCtx *common.CommonContext) ([]viewobject.ReviewVO, error)
	CourseReviews(commonCtx *common.CommonContext, courseID int) ([]viewobject.ReviewVO, error)
}

type reviewQueryService struct {
	reviewRepo review.ReviewRepository
	courseRepo review.CourseRepository
}

func NewReviewQueryService(
	reviewRepo review.ReviewRepository,
	courseRepo review.CourseRepository,
) ReviewQueryService {
	return &reviewQueryService{
		reviewRepo: reviewRepo,
		courseRepo: courseRepo,
	}
}

func (s *reviewQueryService) LatestReviews(commonCtx *common.CommonContext) ([]viewobject.ReviewVO, error) {
	reviews, err := s.reviewRepo.FindBy(commonCtx.Ctx, review.ReviewFilter{})
	if err != nil {
		return nil, apperror.ErrDB
	}
	return s.listReviews(commonCtx, reviews, true), nil
}

func (s *reviewQueryService) CourseReviews(commonCtx *common.CommonContext, courseID int) ([]viewobject.ReviewVO, error) {
	reviews, err := s.reviewRepo.FindBy(commonCtx.Ctx, review.ReviewFilter{CourseID: &courseID})
	if err != nil {
		return nil, apperror.ErrDB
	}
	return s.listReviews(commonCtx, reviews, false), nil
}

func (s *reviewQueryService) listReviews(commonCtx *common.CommonContext, reviews []review.Review, withCourse bool) []viewobject.ReviewVO {
	reviewList := make([]viewobject.ReviewVO, len(reviews))
	for i, r := range reviews {
		reviewList[i] = viewobject.NewReviewVO(&r, withCourse)
	}
	return reviewList
}
