package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/auth"
	authcommand "jcourse_go/internal/application/auth/command"
	authquery "jcourse_go/internal/application/auth/query"
	domainauth "jcourse_go/internal/domain/auth"
)

type AuthController struct {
	authCommandService authcommand.AuthCommandService
	authQueryService   authquery.AuthQueryService
	codeService        auth.VerificationCodeService
}

func NewAuthController(
	authCommandService authcommand.AuthCommandService,
	authQueryService authquery.AuthQueryService,
	codeService auth.VerificationCodeService,
) *AuthController {
	return &AuthController{
		authCommandService: authCommandService,
		authQueryService:   authQueryService,
		codeService:        codeService,
	}
}

type AuthResponse struct {
	SessionID string `json:"session_id"`
}

func NewAuthResponse(sessionID string) AuthResponse {
	return AuthResponse{SessionID: sessionID}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var cmd domainauth.LoginCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	err := c.authCommandService.Login(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	// Get session from context (created by auth command service)
	sessionID := ctx.GetString("session_id")
	response := NewAuthResponse(sessionID)
	HandleSuccess(ctx, response)
}

func (c *AuthController) Register(ctx *gin.Context) {
	var cmd domainauth.RegisterCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	err := c.authCommandService.Register(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	// Get session from context (created by auth command service)
	sessionID := ctx.GetString("session_id")
	response := NewAuthResponse(sessionID)
	HandleSuccessWithStatus(ctx, http.StatusCreated, response)
}

func (c *AuthController) Logout(ctx *gin.Context) {
	sessionID := ctx.GetHeader("X-Session-ID")
	if sessionID == "" {
		HandleValidationError(ctx, "session_id required")
		return
	}

	cmd := domainauth.LogoutCommand{SessionID: sessionID}
	err := c.authCommandService.Logout(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}

func (c *AuthController) SendVerificationCode(ctx *gin.Context) {
	var cmd domainauth.SendVerificationCodeCommand
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		HandleValidationError(ctx, "invalid request body")
		return
	}

	err := c.authCommandService.SendVerificationCode(ctx, cmd)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, nil)
}
