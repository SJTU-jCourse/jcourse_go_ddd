package auth

import (
	"time"

	"jcourse_go/internal/domain/common"
)

type User struct {
	ID       int
	Username string
	Password string
	Email    string

	Role        common.Role
	LastSeenAt  time.Time
	SuspendedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (u *User) IsSuspended() bool {
	return u.SuspendedAt != nil
}

func (u *User) IsAdmin() bool {
	return u.Role == common.RoleAdmin
}

type VerificationCode struct {
	Code      string
	Email     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

func (c *VerificationCode) IsExpired(now time.Time) bool {
	return now.After(c.ExpiresAt)
}

func (c *VerificationCode) Validate(code string) bool {
	return c.Code == code && !c.IsExpired(time.Now())
}
