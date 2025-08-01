package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/point"
	"jcourse_go/internal/infrastructure/entity"
)

type userPointRepository struct {
	db *gorm.DB
}

func NewUserPointRepository(db *gorm.DB) point.UserPointRepository {
	return &userPointRepository{db: db}
}

func (r *userPointRepository) GetUserAllPoints(ctx context.Context, userID int) (*point.UserPoint, error) {
	var recordEntities []entity.UserPointRecord
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&recordEntities)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user point records: %w", result.Error)
	}

	records := make([]point.UserPointRecord, len(recordEntities))
	totalPoint := 0
	for i, recordEntity := range recordEntities {
		records[i] = *r.toDomainPointRecord(&recordEntity)
		totalPoint += recordEntity.Point
	}

	return &point.UserPoint{
		TotalPoint: totalPoint,
		Records:    records,
	}, nil
}

func (r *userPointRepository) GetPointRecord(ctx context.Context, itemID int) (*point.UserPointRecord, error) {
	var recordEntity entity.UserPointRecord
	result := r.db.WithContext(ctx).First(&recordEntity, itemID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get point record: %w", result.Error)
	}
	return r.toDomainPointRecord(&recordEntity), nil
}

func (r *userPointRepository) Save(ctx context.Context, point *point.UserPointRecord) error {
	pointEntity := r.toORMPointRecord(point)
	result := r.db.WithContext(ctx).Create(pointEntity)
	if result.Error != nil {
		return fmt.Errorf("failed to save point record: %w", result.Error)
	}
	return nil
}

// Helper methods to convert between domain and ORM models
func (r *userPointRepository) toDomainPointRecord(recordEntity *entity.UserPointRecord) *point.UserPointRecord {
	return &point.UserPointRecord{
		ItemID:      recordEntity.ID,
		UserID:      recordEntity.UserID,
		Point:       recordEntity.Point,
		Description: recordEntity.Description,
		CreatedAt:   recordEntity.CreatedAt,
	}
}

func (r *userPointRepository) toORMPointRecord(point *point.UserPointRecord) *entity.UserPointRecord {
	return &entity.UserPointRecord{
		ID:          point.ItemID,
		UserID:      point.UserID,
		Point:       point.Point,
		Action:      "point_action", // Default action
		Description: point.Description,
		CreatedAt:   point.CreatedAt,
	}
}
