package announcement

import (
	"context"
	"time"
)

type Announcement struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"is_published"`
	PublishedAt time.Time `json:"published_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AnnouncementRepository interface {
	FindPublished(ctx context.Context) ([]Announcement, error)
}
