package repository

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/announcement"
	"jcourse_go/internal/infrastructure/entity"
)

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) announcement.AnnouncementRepository {
	return &announcementRepository{db: db}
}

func (r *announcementRepository) FindPublished(ctx context.Context) ([]announcement.Announcement, error) {
	var announcementEntities []entity.Announcement
	result := r.db.Where("is_published = ?", true).Order("published_at DESC").Find(&announcementEntities)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return []announcement.Announcement{}, nil
		}
		return nil, fmt.Errorf("failed to get published announcements: %w", result.Error)
	}

	announcements := make([]announcement.Announcement, len(announcementEntities))
	for i, entity := range announcementEntities {
		announcements[i] = announcement.Announcement{
			ID:          entity.ID,
			Title:       entity.Title,
			Content:     entity.Content,
			IsPublished: entity.IsPublished,
			PublishedAt: entity.PublishedAt,
			CreatedAt:   entity.CreatedAt,
			UpdatedAt:   entity.UpdatedAt,
		}
	}

	return announcements, nil
}
