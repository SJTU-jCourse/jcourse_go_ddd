package query

import (
	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
	"jcourse_go/pkg/apperror"
)

type CourseQueryService interface {
	FindCoursesBy(commonCtx *common.CommonContext, filter review.CourseFilter) ([]viewobject.CourseListItemVO, error)
	GetCourseDetail(commonCtx *common.CommonContext, courseID int) (*viewobject.CourseDetailVO, error)
}

type courseQueryService struct {
	courseRepo review.CourseRepository
	reviewRepo review.ReviewRepository
}

func NewCourseQueryService(
	courseRepo review.CourseRepository,
	reviewRepo review.ReviewRepository) CourseQueryService {
	return &courseQueryService{
		courseRepo: courseRepo,
		reviewRepo: reviewRepo,
	}
}

func (s *courseQueryService) FindCoursesBy(commonCtx *common.CommonContext, filter review.CourseFilter) ([]viewobject.CourseListItemVO, error) {
	courses, err := s.courseRepo.FindBy(commonCtx.Ctx, filter)
	if err != nil {
		return nil, apperror.ErrDB
	}
	courseList := make([]viewobject.CourseListItemVO, len(courses))
	for i, c := range courses {
		courseList[i] = viewobject.NewCourseListItemVO(&c)
	}
	return courseList, nil
}

func (s *courseQueryService) GetCourseDetail(commonCtx *common.CommonContext, courseID int) (*viewobject.CourseDetailVO, error) {
	course, err := s.courseRepo.Get(commonCtx.Ctx, courseID)
	if err != nil {
		return nil, apperror.ErrDB
	}
	courseDetailVO := viewobject.NewCourseDetailVO(course)
	return &courseDetailVO, nil
}
