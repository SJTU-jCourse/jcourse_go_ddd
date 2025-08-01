package command

import (
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
)

type CourseCommandService interface {
	AddUserEnrolledCourse(commonCtx *common.CommonContext, courseID int) error
	WatchCourse(commonCtx *common.CommonContext, courseID int, watch bool) error
}

type courseCommandService struct {
	courseRepo review.CourseRepository
}

func NewCourseCommandService(
	courseRepo review.CourseRepository,
) CourseCommandService {
	return &courseCommandService{
		courseRepo: courseRepo,
	}
}

func (s *courseCommandService) AddUserEnrolledCourse(commonCtx *common.CommonContext, courseID int) error {
	return s.courseRepo.AddUserEnrolledCourse(commonCtx.Ctx, commonCtx.User.UserID, courseID)
}

func (s *courseCommandService) WatchCourse(commonCtx *common.CommonContext, courseID int, watch bool) error {
	return s.courseRepo.WatchCourse(commonCtx.Ctx, commonCtx.User.UserID, courseID, watch)
}
