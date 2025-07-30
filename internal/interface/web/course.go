package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/review/query"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
)

type CourseController struct {
	courseQueryService query.CourseQueryService
}

func NewCourseController(courseQueryService query.CourseQueryService) *CourseController {
	return &CourseController{
		courseQueryService: courseQueryService,
	}
}

func (c *CourseController) GetCourseDetail(ctx *gin.Context) {
	courseIDStr := ctx.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid course id")
		return
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 0, Role: common.RoleUser},
	}

	courseDetail, err := c.courseQueryService.GetCourseDetail(commonCtx, courseID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, courseDetail)
}

func (c *CourseController) SearchCourses(ctx *gin.Context) {
	var filter review.CourseFilter

	if name := ctx.Query("name"); name != "" {
		filter.Name = &name
	}

	if code := ctx.Query("code"); code != "" {
		filter.Code = &code
	}

	if dept := ctx.Query("department"); dept != "" {
		filter.Departments = []string{dept}
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 0, Role: common.RoleUser},
	}

	courses, err := c.courseQueryService.FindCoursesBy(commonCtx, filter)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, gin.H{"courses": courses})
}
