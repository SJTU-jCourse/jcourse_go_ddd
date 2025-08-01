package web

import (
	"github.com/gin-gonic/gin"

	authquery "jcourse_go/internal/application/auth/query"
	reviewquery "jcourse_go/internal/application/review/query"
	"jcourse_go/internal/interface/dto"
)

type UserController struct {
	userQueryService   authquery.UserQueryService
	reviewQueryService reviewquery.ReviewQueryService
}

func NewUserController(userQueryService authquery.UserQueryService, reviewQueryService reviewquery.ReviewQueryService) *UserController {
	return &UserController{
		userQueryService:   userQueryService,
		reviewQueryService: reviewQueryService,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	userInfo, err := c.userQueryService.GetUserInfo(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, userInfo)
}

func (c *UserController) UpdateUserInfo(ctx *gin.Context) {
	var cmd dto.UpdateUserInfoRequest

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := GetCommonContext(ctx)

	err := c.userQueryService.UpdateUserInfo(commonCtx, cmd.Nickname)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *UserController) GetUserReviews(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	reviews, err := c.reviewQueryService.GetUserReviews(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, reviews)
}
