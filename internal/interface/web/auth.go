package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appauth "jcourse_go/internal/application/auth"
	"jcourse_go/internal/domain/auth"
)

type AuthController struct {
	authService appauth.AuthService
	codeService appauth.VerificationCodeService
}

func NewAuthController(authService appauth.AuthService, codeService appauth.VerificationCodeService) *AuthController {
	return &AuthController{
		authService: authService,
		codeService: codeService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var cmd auth.LoginCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	sessionID, err := c.authService.Login(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, gin.H{"session_id": *sessionID})
}

func (c *AuthController) Register(ctx *gin.Context) {
	var cmd auth.RegisterCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	sessionID, err := c.authService.Register(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccessWithStatus(ctx, http.StatusCreated, gin.H{"session_id": *sessionID})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	sessionID := ctx.GetHeader("X-Session-ID")
	if sessionID == "" {
		HandleValidationError(ctx, "session_id required")
		return
	}

	cmd := auth.LogoutCommand{SessionID: sessionID}
	err := c.authService.Logout(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *AuthController) SendVerificationCode(ctx *gin.Context) {
	var cmd auth.SendVerificationCodeCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	err := c.authService.SendVerificationCode(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}
