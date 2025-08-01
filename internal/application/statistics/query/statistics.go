package query

import (
	"time"

	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/statistics"
	"jcourse_go/pkg/apperror"
)

type StatisticsQueryService interface {
	GetCurrentStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error)
	GetDailyStatistics(commonCtx *common.CommonContext, date time.Time) (*viewobject.DailyStatisticsVO, error)
	GetDailyStatisticsRange(commonCtx *common.CommonContext, startDate, endDate time.Time) ([]viewobject.DailyStatisticsVO, error)
	GetLatestDailyStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error)
}

type statisticsQueryService struct {
	statisticsRepo statistics.StatisticsRepository
}

func NewStatisticsQueryService(statisticsRepo statistics.StatisticsRepository) StatisticsQueryService {
	return &statisticsQueryService{
		statisticsRepo: statisticsRepo,
	}
}

func (s *statisticsQueryService) GetCurrentStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetCurrentStatistics(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	statisticsVO := viewobject.NewDailyStatisticsVO(stats)
	return &statisticsVO, nil
}

func (s *statisticsQueryService) GetDailyStatistics(commonCtx *common.CommonContext, date time.Time) (*viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetDailyStatistics(commonCtx.Ctx, date)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	statisticsVO := viewobject.NewDailyStatisticsVO(stats)
	return &statisticsVO, nil
}

func (s *statisticsQueryService) GetDailyStatisticsRange(commonCtx *common.CommonContext, startDate, endDate time.Time) ([]viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetDailyStatisticsRange(commonCtx.Ctx, startDate, endDate)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	return viewobject.NewDailyStatisticsVOList(stats), nil
}

func (s *statisticsQueryService) GetLatestDailyStatistics(commonCtx *common.CommonContext) (*viewobject.DailyStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetLatestDailyStatistics(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	statisticsVO := viewobject.NewDailyStatisticsVO(stats)
	return &statisticsVO, nil
}
