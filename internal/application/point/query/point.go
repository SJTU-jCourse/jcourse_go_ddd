package query

import (
	"context"

	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/point"
)

type UserPointQueryService interface {
	GetUserPoint(ctx context.Context, userID int) (*viewobject.UserPointVO, error)
}

func NewUserPointQueryService(pointRepo point.UserPointRepository) UserPointQueryService {
	return &userPointQueryService{
		pointRepo: pointRepo,
	}
}

type userPointQueryService struct {
	pointRepo point.UserPointRepository
}

func (s *userPointQueryService) GetUserPoint(ctx context.Context, userID int) (*viewobject.UserPointVO, error) {
	userPoint, err := s.pointRepo.GetUserAllPoints(ctx, userID)
	if err != nil {
		return nil, err
	}
	vo := viewobject.NewUserPointVO(userPoint)
	return &vo, nil
}
