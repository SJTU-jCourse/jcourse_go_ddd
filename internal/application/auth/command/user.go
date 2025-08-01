package command

import (
	"context"
	"time"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/pkg/apperror"
)

type UserCommandService interface {
	UpdateUserInfo(ctx context.Context, userID int, nickname string) error
}

type userCommandService struct {
	userRepo auth.UserRepository
}

func NewUserCommandService(
	userRepo auth.UserRepository,
) UserCommandService {
	return &userCommandService{
		userRepo: userRepo,
	}
}

func (s *userCommandService) UpdateUserInfo(ctx context.Context, userID int, nickname string) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "update_user_info").WithMetadata("user_id", userID)
	}
	if user == nil {
		return apperror.ErrUserNotFound.WithMetadata("user_id", userID)
	}

	user.UpdateNickname(nickname)
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}
