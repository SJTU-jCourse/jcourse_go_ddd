package web

import (
	"github.com/gin-gonic/gin"

	authcommand "jcourse_go/internal/application/auth/command"
	authquery "jcourse_go/internal/application/auth/query"
	reviewquery "jcourse_go/internal/application/review/query"
	"jcourse_go/internal/interface/dto"
)

type UserController struct {
	userCommandService authcommand.UserCommandService
	userQueryService   authquery.UserQueryService
	reviewQueryService reviewquery.ReviewQueryService
}

func NewUserController(
	userCommandService authcommand.UserCommandService,
	userQueryService authquery.UserQueryService,
	reviewQueryService reviewquery.ReviewQueryService,
) *UserController {
	return &UserController{
		userCommandService: userCommandService,
		userQueryService:   userQueryService,
		reviewQueryService: reviewQueryService,
	}
}

func (c *UserController) GetUserInfo(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	userInfo, err := c.userQueryService.GetUserInfo(ctx, commonCtx.User.UserID)
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

	err := c.userCommandService.UpdateUserInfo(ctx, commonCtx.User.UserID, cmd.Nickname)
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
