package service

import (
	"context"
	"time"

	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/statistics"
	"jcourse_go/pkg/apperror"
)

type DailyStatisticsService interface {
	GetDailyStatistics(commonCtx *common.CommonContext, date time.Time) (*viewobject.DailyStatisticsVO, error)
	GetDailyStatisticsRange(commonCtx *common.CommonContext, startDate, endDate time.Time) ([]viewobject.DailyStatisticsVO, error)
	GetLatestDailyStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error)
	CalculateAndSaveDailyStatistics(ctx context.Context, date time.Time) error
}

type dailyStatisticsService struct {
	statisticsRepo statistics.StatisticsRepository
}

func NewDailyStatisticsService(statisticsRepo statistics.StatisticsRepository) DailyStatisticsService {
	return &dailyStatisticsService{
		statisticsRepo: statisticsRepo,
	}
}

func (s *dailyStatisticsService) GetDailyStatistics(commonCtx *common.CommonContext, date time.Time) (*viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetDailyStatistics(commonCtx.Ctx, date)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	vo := viewobject.NewDailyStatisticsVO(stats)
	return &vo, nil
}

func (s *dailyStatisticsService) GetDailyStatisticsRange(commonCtx *common.CommonContext, startDate, endDate time.Time) ([]viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetDailyStatisticsRange(commonCtx.Ctx, startDate, endDate)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	return viewobject.NewDailyStatisticsVOList(stats), nil
}

func (s *dailyStatisticsService) GetLatestDailyStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetLatestDailyStatistics(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	vo := viewobject.NewDailyStatisticsVO(stats)
	return &vo, nil
}

func (s *dailyStatisticsService) CalculateAndSaveDailyStatistics(ctx context.Context, date time.Time) error {
	// Get current statistics for the given date (this will calculate real-time if not available)
	currentStats, err := s.statisticsRepo.GetCurrentStatistics(ctx)
	if err != nil {
		return apperror.ErrDB.Wrap(err)
	}

	// Create daily statistics with the calculated data
	dailyStats := &statistics.DailyStatistics{
		Date:               date,
		DAU:                currentStats.DAU,
		DNU:                currentStats.DNU,
		MAU:                currentStats.MAU,
		DailyNewReviews:    currentStats.DailyNewReviews,
		TotalReviews:       currentStats.TotalReviews,
		TotalCourses:       currentStats.TotalCourses,
		CoursesWithReviews: currentStats.CoursesWithReviews,
	}

	// Save daily statistics
	err = s.statisticsRepo.SaveDailyStatistics(ctx, dailyStats)
	if err != nil {
		return apperror.ErrDB.Wrap(err)
	}

	return nil
}
