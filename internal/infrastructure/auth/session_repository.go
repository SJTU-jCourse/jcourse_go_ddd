package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"

	"jcourse_go/internal/domain/auth"
)

type sessionRepository struct {
	redis *redis.Client
}

func NewSessionRepository(redisClient *redis.Client) auth.SessionRepository {
	return &sessionRepository{redis: redisClient}
}

func (r *sessionRepository) Store(ctx context.Context, userID int) (string, error) {
	sessionID := fmt.Sprintf("session:%d:%d", userID, time.Now().Unix())

	err := r.redis.Set(ctx, sessionID, userID, 24*time.Hour).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session: %w", err)
	}

	return sessionID, nil
}

func (r *sessionRepository) Get(ctx context.Context, sessionID string) (int, error) {
	userIDStr, err := r.redis.Get(ctx, sessionID).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	var userID int
	_, err = fmt.Sscanf(userIDStr, "%d", &userID)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user id: %w", err)
	}

	return userID, nil
}

func (r *sessionRepository) Delete(ctx context.Context, sessionID string) error {
	err := r.redis.Del(ctx, sessionID).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}
