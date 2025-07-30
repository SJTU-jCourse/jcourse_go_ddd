package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/point/command"
	"jcourse_go/internal/application/point/query"
	"jcourse_go/internal/domain/common"
)

type UserPointController struct {
	pointCommandService command.PointCommandService
	pointQueryService   query.UserPointQueryService
}

func NewUserPointController(pointCommandService command.PointCommandService, pointQueryService query.UserPointQueryService) *UserPointController {
	return &UserPointController{
		pointCommandService: pointCommandService,
		pointQueryService:   pointQueryService,
	}
}

func (c *UserPointController) GetUserPoint(ctx *gin.Context) {
	userIDStr := ctx.Param("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		HandleValidationError(ctx, "invalid user id")
		return
	}

	commonCtx := &common.CommonContext{
		Ctx:  ctx,
		User: &common.User{UserID: userID, Role: common.RoleUser},
	}

	userPoint, err := c.pointQueryService.GetUserPoint(commonCtx.Ctx, userID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, userPoint)
}
