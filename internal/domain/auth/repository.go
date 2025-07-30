package auth

import "context"

type CodeRepository interface {
	Get(ctx context.Context, email string) (*VerificationCode, error)
	Save(ctx context.Context, code *VerificationCode) error
	Delete(ctx context.Context, code *VerificationCode) error
}

type UserFilter struct {
	UserIDs []int
}

type UserRepository interface {
	Get(ctx context.Context, email string) (*User, error)
	FindBy(ctx context.Context, filter UserFilter) ([]User, error)
	Save(ctx context.Context, user *User) (int, error)
}

type SessionRepository interface {
	Store(ctx context.Context, userID int) (string, error)
	Get(ctx context.Context, sessionID string) (int, error)
	Delete(ctx context.Context, sessionID string) error
}
