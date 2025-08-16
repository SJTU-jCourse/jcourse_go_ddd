package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/infrastructure/entity"
)

const (
	SessionExpiration = 24 * time.Hour
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) auth.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Store(ctx context.Context, userID int) (string, error) {
	sessionID := fmt.Sprintf("session:%d:%d", userID, time.Now().Unix())
	expiresAt := time.Now().Add(SessionExpiration)

	session := entity.UserSession{
		UserID:    userID,
		Token:     sessionID,
		ExpiresAt: expiresAt,
	}

	if err := r.db.Create(&session).Error; err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return sessionID, nil
}

func (r *sessionRepository) Get(ctx context.Context, sessionID string) (int, error) {
	var session entity.UserSession
	err := r.db.Where("token = ? AND expires_at > ?", sessionID, time.Now()).First(&session).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("session not found or expired")
		}
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	return session.UserID, nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	err := r.db.Where("token = ?", sessionID).Delete(&entity.UserSession{}).Error
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
