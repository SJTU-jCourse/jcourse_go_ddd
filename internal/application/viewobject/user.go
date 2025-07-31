package viewobject

import "jcourse_go/internal/domain/auth"

type UserInfoVO struct {
	UserID   int    `json:"user_id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
}

func NewUserInfoVO(user *auth.User) UserInfoVO {
	return UserInfoVO{
		UserID:   user.ID,
		Email:    user.Email,
		Nickname: user.Username, // Using Username as Nickname
		Role:     string(user.Role),
	}
}
