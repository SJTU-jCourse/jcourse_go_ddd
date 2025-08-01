package web

import (
	"time"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/statistics/query"
	"jcourse_go/internal/application/statistics/service"
)

type StatisticsController struct {
	statisticsQueryService query.StatisticsQueryService
	dailyStatisticsService service.DailyStatisticsService
}

func NewStatisticsController(statisticsQueryService query.StatisticsQueryService, dailyStatisticsService service.DailyStatisticsService) *StatisticsController {
	return &StatisticsController{
		statisticsQueryService: statisticsQueryService,
		dailyStatisticsService: dailyStatisticsService,
	}
}

func (c *StatisticsController) GetSystemStatistics(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	statistics, err := c.statisticsQueryService.GetCurrentStatistics(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, statistics)
}

func (c *StatisticsController) GetDailyStatistics(ctx *gin.Context) {
	dateStr := ctx.Param("date")
	if dateStr == "" {
		HandleValidationError(ctx, "date parameter is required")
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		HandleValidationError(ctx, "invalid date format, use YYYY-MM-DD")
		return
	}

	commonCtx := GetCommonContext(ctx)
	statistics, err := c.statisticsQueryService.GetDailyStatistics(commonCtx, date)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, statistics)
}

func (c *StatisticsController) GetDailyStatisticsRange(ctx *gin.Context) {
	startDateStr := ctx.Query("start_date")
	endDateStr := ctx.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		HandleValidationError(ctx, "start_date and end_date parameters are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		HandleValidationError(ctx, "invalid start_date format, use YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		HandleValidationError(ctx, "invalid end_date format, use YYYY-MM-DD")
		return
	}

	// Limit the range to prevent excessive data retrieval
	if int(endDate.Sub(startDate).Hours()/24) > 365 {
		HandleValidationError(ctx, "date range cannot exceed 365 days")
		return
	}

	commonCtx := GetCommonContext(ctx)
	statistics, err := c.statisticsQueryService.GetDailyStatisticsRange(commonCtx, startDate, endDate)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, statistics)
}

func (c *StatisticsController) GetLatestDailyStatistics(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)
	statistics, err := c.statisticsQueryService.GetLatestDailyStatistics(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, statistics)
}

func (c *StatisticsController) TriggerDailyStatistics(ctx *gin.Context) {
	dateStr := ctx.Query("date")
	var date time.Time
	var err error

	if dateStr == "" {
		// Default to yesterday
		date = time.Now().AddDate(0, 0, -1)
		date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			HandleValidationError(ctx, "invalid date format, use YYYY-MM-DD")
			return
		}
	}

	commonCtx := GetCommonContext(ctx)
	err = c.dailyStatisticsService.CalculateAndSaveDailyStatistics(commonCtx.Ctx, date)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, gin.H{
		"message": "Daily statistics calculated successfully",
		"date":    date.Format("2006-01-02"),
	})
}
