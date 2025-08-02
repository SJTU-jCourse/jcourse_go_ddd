package viewobject

import (
	"time"

	"jcourse_go/internal/domain/common"
)

type UserInfoVO struct {
	UserID      int         `json:"user_id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	Role        common.Role `json:"role"`
	IsSuspended bool        `json:"is_suspended"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
