package common

import (
	"context"
)

type CommonContext struct {
	Ctx  context.Context
	User *User
}

type User struct {
	UserID int
	Role   Role
}

func NewCommonContext(ctx context.Context, user *User) *CommonContext {
	return &CommonContext{
		Ctx:  ctx,
		User: user,
	}
}
