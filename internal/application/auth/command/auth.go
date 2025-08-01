package command

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/auth"
	domainauth "jcourse_go/internal/domain/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
	"jcourse_go/pkg/password"
)

type AuthCommandService interface {
	Login(ctx context.Context, cmd domainauth.LoginCommand) error
	Register(ctx context.Context, cmd domainauth.RegisterCommand) error
	Logout(ctx context.Context, cmd domainauth.LogoutCommand) error
	SendVerificationCode(ctx context.Context, cmd domainauth.SendVerificationCodeCommand) error
}

func NewAuthCommandService(
	userRepo domainauth.UserRepository,
	hasher password.Hasher,
	session domainauth.SessionRepository,
	codeService auth.VerificationCodeService,
) AuthCommandService {
	return &authCommandService{
		userRepo:    userRepo,
		hasher:      hasher,
		session:     session,
		codeService: codeService,
	}
}

type authCommandService struct {
	userRepo    domainauth.UserRepository
	hasher      password.Hasher
	session     domainauth.SessionRepository
	codeService auth.VerificationCodeService
}

func (s *authCommandService) Login(ctx context.Context, cmd domainauth.LoginCommand) error {
	user, err := s.userRepo.Get(ctx, cmd.Email)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "login").WithMetadata("email", cmd.Email)
	}
	if user == nil {
		return apperror.ErrWrongAuth.WithUserMessage("Invalid email or password").WithMetadata("email", cmd.Email)
	}
	if err := s.hasher.Validate(cmd.Password, user.Password); err != nil {
		return apperror.ErrWrongAuth.WithUserMessage("Invalid email or password").WithMetadata("email", cmd.Email)
	}
	if user.IsSuspended() {
		return apperror.ErrSuspended.WithMetadata("user_id", user.ID)
	}
	sessionID, err := s.session.Store(ctx, user.ID)
	if err != nil {
		return apperror.ErrSession.Wrap(err).WithMetadata("operation", "login").WithMetadata("user_id", user.ID)
	}
	// Store session ID in context for controller to access
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set("session_id", sessionID)
	}
	return nil
}

func (s *authCommandService) Logout(ctx context.Context, cmd domainauth.LogoutCommand) error {
	if err := s.session.Delete(ctx, cmd.SessionID); err != nil {
		return apperror.ErrSession.Wrap(err).WithMetadata("operation", "logout").WithMetadata("session_id", cmd.SessionID)
	}
	return nil
}

func (s *authCommandService) newUserFromRegister(cmd domainauth.RegisterCommand) *domainauth.User {
	now := time.Now()
	hashedPassword := s.hasher.Hash(cmd.Password)
	user := &domainauth.User{
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

func (s *authCommandService) Register(ctx context.Context, cmd domainauth.RegisterCommand) error {
	if err := s.codeService.Verify(ctx, cmd.Code, cmd.Email); err != nil {
		return apperror.ErrWrongAuth.Wrap(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	existUser, err := s.userRepo.Get(ctx, cmd.Email)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	if existUser != nil {
		return apperror.ErrWrongAuth.WithUserMessage("Email already registered").WithMetadata("email", cmd.Email)
	}
	user := s.newUserFromRegister(cmd)
	userID, err := s.userRepo.Save(ctx, user)
	if err != nil {
		return apperror.WrapDB(err).WithMetadata("operation", "register").WithMetadata("email", cmd.Email)
	}
	sessionID, err := s.session.Store(ctx, userID)
	if err != nil {
		return apperror.ErrSession.Wrap(err).WithMetadata("operation", "register").WithMetadata("user_id", userID)
	}
	// Store session ID in context for controller to access
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ginCtx.Set("session_id", sessionID)
	}
	return nil
}

func (s *authCommandService) SendVerificationCode(ctx context.Context, cmd domainauth.SendVerificationCodeCommand) error {
	return s.codeService.SendCode(ctx, cmd.Email)
}
