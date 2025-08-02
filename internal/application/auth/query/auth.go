package query

import (
	"context"

	"jcourse_go/internal/application/viewobject"
	domainauth "jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
)

type AuthQueryService interface {
	GetUserFromSession(ctx context.Context, sessionID string) (*common.User, error)
	GetUserInfo(ctx context.Context, userID int) (*viewobject.UserInfoVO, error)
}

type authQueryService struct {
	userRepo domainauth.UserRepository
	session  domainauth.SessionRepository
}

func NewAuthQueryService(
	userRepo domainauth.UserRepository,
	session domainauth.SessionRepository,
) AuthQueryService {
	return &authQueryService{
		userRepo: userRepo,
		session:  session,
	}
}

func (s *authQueryService) GetUserFromSession(ctx context.Context, sessionID string) (*common.User, error) {
	userID, err := s.session.Get(ctx, sessionID)
	if err != nil {
		return nil, apperror.ErrSession.Wrap(err).WithMetadata("operation", "get_user_from_session").WithMetadata("session_id", sessionID)
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.WrapDB(err).WithMetadata("operation", "get_user_from_session").WithMetadata("user_id", userID)
	}
	if user == nil {
		return nil, apperror.ErrUserNotFound.WithMetadata("user_id", userID)
	}

	return &common.User{
		UserID: user.ID,
		Role:   user.Role,
	}, nil
}

func (s *authQueryService) GetUserInfo(ctx context.Context, userID int) (*viewobject.UserInfoVO, error) {
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
