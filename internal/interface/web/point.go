package web

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/point/command"
	"jcourse_go/internal/application/point/query"
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
	var cmd struct {
		UserID int    `json:"user_id" binding:"required"`
		Amount int    `json:"amount" binding:"required"`
		Reason string `json:"reason" binding:"required"`
	}

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
	var cmd struct {
		FromUserID int    `json:"from_user_id" binding:"required"`
		ToUserID   int    `json:"to_user_id" binding:"required"`
		Amount     int    `json:"amount" binding:"required"`
		Reason     string `json:"reason" binding:"required"`
	}

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
