package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/review/query"
	"jcourse_go/internal/application/viewobject"
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

	commonCtx := GetCommonContext(ctx)

	courseDetail, err := c.courseQueryService.GetCourseDetail(commonCtx, courseID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, courseDetail)
}

type CourseResponse struct {
	Courses []viewobject.CourseListItemVO `json:"courses"`
}

func NewCourseResponse(courses []viewobject.CourseListItemVO) CourseResponse {
	return CourseResponse{Courses: courses}
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

	commonCtx := GetCommonContext(ctx)

	courses, err := c.courseQueryService.FindCoursesBy(commonCtx, filter)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	response := NewCourseResponse(courses)
	HandleSuccess(ctx, response)
}

func (c *CourseController) GetCourseFilter(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	filter, err := c.courseQueryService.GetCourseFilter(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, filter)
}

func (c *CourseController) GetUserEnrolledCourses(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	courses, err := c.courseQueryService.GetUserEnrolledCourses(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	response := NewCourseResponse(courses)
	HandleSuccess(ctx, response)
}

func (c *CourseController) AddUserEnrolledCourse(ctx *gin.Context) {
	var cmd struct {
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := GetCommonContext(ctx)

	err := c.courseQueryService.AddUserEnrolledCourse(commonCtx, cmd.CourseID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *CourseController) WatchCourse(ctx *gin.Context) {
	courseIDStr := ctx.Param("id")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid course id")
		return
	}

	var cmd struct {
		Watch bool `json:"watch" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := GetCommonContext(ctx)

	err = c.courseQueryService.WatchCourse(commonCtx, courseID, cmd.Watch)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}
