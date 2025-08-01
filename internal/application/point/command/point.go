package command

import (
	"context"

	"jcourse_go/internal/domain/point"
)

type PointCommandService interface {
	CreatePoint(ctx context.Context) error
	Transaction(ctx context.Context) error
}

func NewPointCommandService(repo point.UserPointRepository) PointCommandService {
	return &pointCommandService{
		repo: repo,
	}
}

type pointCommandService struct {
	repo point.UserPointRepository
}

func (s *pointCommandService) CreatePoint(ctx context.Context) error {
	return nil
}

func (s *pointCommandService) Transaction(ctx context.Context) error {
	return nil
}
