package auth

import (
	"context"
	"time"

	"jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
	"jcourse_go/pkg/password"
)

type AuthService interface {
	Login(ctx context.Context, cmd auth.LoginCommand) (*string, error)
	Register(ctx context.Context, cmd auth.RegisterCommand) (*string, error)
	Logout(ctx context.Context, cmd auth.LogoutCommand) error
	SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error
	GetUserFromSession(ctx context.Context, sessionID string) (*common.User, error)
}

func NewAuthService(
	userRepo auth.UserRepository,
	hasher password.Hasher,
	session auth.SessionRepository,
	codeService VerificationCodeService,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		hasher:      hasher,
		session:     session,
		codeService: codeService,
	}
}

type authService struct {
	userRepo    auth.UserRepository
	hasher      password.Hasher
	session     auth.SessionRepository
	codeService VerificationCodeService
}

func (s *authService) Login(ctx context.Context, cmd auth.LoginCommand) (*string, error) {
	user, err := s.userRepo.Get(ctx, cmd.Email)
	if err != nil {
		return nil, apperror.ErrWrongAuth
	}
	if err := s.hasher.Validate(cmd.Password, user.Password); err != nil {
		return nil, apperror.ErrWrongAuth
	}
	if user.IsSuspended() {
		return nil, apperror.ErrSuspended
	}
	sessionID, err := s.session.Store(ctx, user.ID)
	if err != nil {
		return nil, apperror.ErrSession
	}
	return &sessionID, nil
}

func (s *authService) Logout(ctx context.Context, cmd auth.LogoutCommand) error {
	if err := s.session.Delete(ctx, cmd.SessionID); err != nil {
		return apperror.ErrSession
	}
	return nil
}

func (s *authService) newUserFromRegister(cmd auth.RegisterCommand) *auth.User {
	now := time.Now()
	hashedPassword := s.hasher.Hash(cmd.Password)
	user := &auth.User{
		Username:   cmd.Email,
		Password:   hashedPassword,
		Email:      cmd.Email,
		Role:       common.RoleUser,
		LastSeenAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	return user
}

func (s *authService) Register(ctx context.Context, cmd auth.RegisterCommand) (*string, error) {
	if err := s.codeService.Verify(ctx, cmd.Code, cmd.Email); err != nil {
		return nil, apperror.ErrWrongAuth
	}
	existUser, err := s.userRepo.Get(ctx, cmd.Email)
	if err != nil {
		return nil, apperror.ErrDB
	}
	if existUser != nil {
		return nil, apperror.ErrWrongAuth
	}
	user := s.newUserFromRegister(cmd)
	userID, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, apperror.ErrDB
	}
	sessionID, err := s.session.Store(ctx, userID)
	if err != nil {
		return nil, apperror.ErrSession
	}
	return &sessionID, nil
}

func (s *authService) SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error {
	return s.codeService.SendCode(ctx, cmd.Email)
}

func (s *authService) GetUserFromSession(ctx context.Context, sessionID string) (*common.User, error) {
	userID, err := s.session.Get(ctx, sessionID)
	if err != nil {
		return nil, apperror.ErrSession
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.ErrUserNotFound
	}

	return &common.User{
		UserID: user.ID,
		Role:   user.Role,
	}, nil
}
