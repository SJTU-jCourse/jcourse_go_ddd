package web

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/statistics/query"
)

type StatisticsController struct {
	statisticsQueryService query.StatisticsQueryService
}

func NewStatisticsController(statisticsQueryService query.StatisticsQueryService) *StatisticsController {
	return &StatisticsController{
		statisticsQueryService: statisticsQueryService,
	}
}

func (c *StatisticsController) GetSystemStatistics(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	statistics, err := c.statisticsQueryService.GetSystemStatistics(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, statistics)
}
