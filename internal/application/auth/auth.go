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
		return nil, apperror.WrapDB(err).WithMetadata("operation", "login").WithMetadata("email", cmd.Email)
	}
	if user == nil {
		return nil, apperror.ErrWrongAuth.WithUserMessage("Invalid email or password").WithMetadata("email", cmd.Email)
	}
	if err := s.hasher.Validate(cmd.Password, user.Password); err != nil {
		return nil, apperror.ErrWrongAuth.WithUserMessage("Invalid email or password").WithMetadata("email", cmd.Email)
	}
	if user.IsSuspended() {
		return nil, apperror.ErrSuspended.WithMetadata("user_id", user.ID)
	}
	sessionID, err := s.session.Store(ctx, user.ID)
	if err != nil {
		return nil, apperror.ErrSession.Wrap(err).WithMetadata("operation", "login").WithMetadata("user_id", user.ID)
	}
	return &sessionID, nil
}

func (s *authService) Logout(ctx context.Context, cmd auth.LogoutCommand) error {
	if err := s.session.Delete(ctx, cmd.SessionID); err != nil {
		return apperror.ErrSession.Wrap(err).WithMetadata("operation", "logout").WithMetadata("session_id", cmd.SessionID)
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
		return nil, apperror.ErrWrongAuth.Wrap(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	existUser, err := s.userRepo.Get(ctx, cmd.Email)
	if err != nil {
		return nil, apperror.WrapDB(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	if existUser != nil {
		return nil, apperror.ErrWrongAuth.WithUserMessage("Email already registered").WithMetadata("email", cmd.Email)
	}
	user := s.newUserFromRegister(cmd)
	userID, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return nil, apperror.WrapDB(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	sessionID, err := s.session.Store(ctx, userID)
	if err != nil {
		return nil, apperror.ErrSession.Wrap(err).WithMetadata("operation", "register").WithMetadata("user_id", userID)
	}
	return &sessionID, nil
}

func (s *authService) SendVerificationCode(ctx context.Context, cmd auth.SendVerificationCodeCommand) error {
	return s.codeService.SendCode(ctx, cmd.Email)
}

func (s *authService) GetUserFromSession(ctx context.Context, sessionID string) (*common.User, error) {
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
