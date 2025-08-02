package query

import (
	"context"

	"jcourse_go/internal/application/viewobject"
	domainauth "jcourse_go/internal/domain/auth"
	"jcourse_go/pkg/apperror"
)

type UserQueryService interface {
	GetUserInfo(ctx context.Context, userID int) (*viewobject.UserInfoVO, error)
}

type userQueryService struct {
	userRepo domainauth.UserRepository
}

func NewUserQueryService(
	userRepo domainauth.UserRepository,
) UserQueryService {
	return &userQueryService{
		userRepo: userRepo,
	}
}

func (s *userQueryService) GetUserInfo(ctx context.Context, userID int) (*viewobject.UserInfoVO, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.WrapDB(err).WithMetadata("operation", "get_user_info").WithMetadata("user_id", userID)
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound.WithMetadata("user_id", userID)
	}

	return &viewobject.UserInfoVO{
		UserID:      user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Role:        user.Role,
		IsSuspended: user.IsSuspended(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}, nil
}
