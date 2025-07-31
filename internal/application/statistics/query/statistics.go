package query

import (
	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/common"
	"jcourse_go/internal/domain/statistics"
	"jcourse_go/pkg/apperror"
)

type StatisticsQueryService interface {
	GetSystemStatistics(commonCtx *common.CommonContext) (*viewobject.SystemStatisticsVO, error)
}

type statisticsQueryService struct {
	statisticsRepo statistics.StatisticsRepository
}

func NewStatisticsQueryService(statisticsRepo statistics.StatisticsRepository) StatisticsQueryService {
	return &statisticsQueryService{
		statisticsRepo: statisticsRepo,
	}
}

func (s *statisticsQueryService) GetSystemStatistics(commonCtx *common.CommonContext) (*viewobject.SystemStatisticsVO, error) {
	stats, err := s.statisticsRepo.GetSystemStatistics(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	statisticsVO := viewobject.NewSystemStatisticsVO(stats)
	return &statisticsVO, nil
}
