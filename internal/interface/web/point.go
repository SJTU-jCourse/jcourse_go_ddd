package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/point/command"
	"jcourse_go/internal/application/point/query"
	"jcourse_go/internal/interface/dto"
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

	commonCtx := GetCommonContext(ctx)

	userPoint, err := c.pointQueryService.GetUserPoint(commonCtx.Ctx, userID)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, userPoint)
}

func (c *UserPointController) CreatePoint(ctx *gin.Context) {
	var cmd dto.CreatePointRequest

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := GetCommonContext(ctx)

	err := c.pointCommandService.CreatePoint(commonCtx, cmd.UserID, cmd.Amount, cmd.Reason)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *UserPointController) Transaction(ctx *gin.Context) {
	var cmd dto.PointTransactionRequest

	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	commonCtx := GetCommonContext(ctx)

	err := c.pointCommandService.Transaction(commonCtx, cmd.FromUserID, cmd.ToUserID, cmd.Amount, cmd.Reason)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}
