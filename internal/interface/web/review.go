package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/review/command"
	"jcourse_go/internal/application/review/query"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/review"
)

type ReviewController struct {
	reviewCommandService command.ReviewCommandService
	reviewQueryService   query.ReviewQueryService
}

func NewReviewController(reviewCommandService command.ReviewCommandService, reviewQueryService query.ReviewQueryService) *ReviewController {
	return &ReviewController{
		reviewCommandService: reviewCommandService,
		reviewQueryService:   reviewQueryService,
	}
}

func (c *ReviewController) WriteReview(ctx *gin.Context) {
	var cmd review.WriteReviewCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 1, Role: common.RoleUser},
	}

	err := c.reviewCommandService.WriteReview(commonCtx, &cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccessWithStatus(ctx, http.StatusCreated, nil)
}

func (c *ReviewController) UpdateReview(ctx *gin.Context) {
	reviewIDStr := ctx.Param("id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid review id")
		return
	}

	var cmd review.UpdateReviewCommand
	cmd.ReviewID = reviewID
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 1, Role: common.RoleUser},
	}

	err = c.reviewCommandService.UpdateReview(commonCtx, &cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *ReviewController) DeleteReview(ctx *gin.Context) {
	reviewIDStr := ctx.Param("id")
	reviewID, err := strconv.Atoi(reviewIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid review id")
		return
	}

	cmd := review.DeleteReviewCommand{ReviewID: reviewID}
	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 1, Role: common.RoleUser},
	}

	err = c.reviewCommandService.DeleteReview(commonCtx, &cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *ReviewController) GetLatestReviews(ctx *gin.Context) {
	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 0, Role: common.RoleUser},
	}

	reviews, err := c.reviewQueryService.LatestReviews(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, gin.H{"reviews": reviews})
}

func (c *ReviewController) GetCourseReviews(ctx *gin.Context) {
	courseIDStr := ctx.Param("courseId")
	courseID, err := strconv.Atoi(courseIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid course id")
		return
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: 0, Role: common.RoleUser},
	}

	reviews, err := c.reviewQueryService.CourseReviews(commonCtx, courseID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, gin.H{"reviews": reviews})
}
