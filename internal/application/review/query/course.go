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
	GetCourseFilter(commonCtx *common.CommonContext) (*viewobject.CourseFilterVO, error)
	GetUserEnrolledCourses(commonCtx *common.CommonContext) ([]viewobject.CourseListItemVO, error)
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

func (s *courseQueryService) GetCourseFilter(commonCtx *common.CommonContext) (*viewobject.CourseFilterVO, error) {
	departments, err := s.courseRepo.GetDepartments(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	categories, err := s.courseRepo.GetCategories(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	filterVO := viewobject.NewCourseFilterVO(departments, categories)
	return &filterVO, nil
}

func (s *courseQueryService) GetUserEnrolledCourses(commonCtx *common.CommonContext) ([]viewobject.CourseListItemVO, error) {
	courseIDs, err := s.courseRepo.GetUserEnrolledCourses(commonCtx.Ctx, commonCtx.User.UserID)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	courseList := make([]viewobject.CourseListItemVO, len(courseIDs))
	for i, courseID := range courseIDs {
		course, err := s.courseRepo.Get(commonCtx.Ctx, courseID)
		if err != nil {
			return nil, apperror.ErrDB.Wrap(err)
		}
		courseList[i] = viewobject.NewCourseListItemVO(course)
	}
	return courseList, nil
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
