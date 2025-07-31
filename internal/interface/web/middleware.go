package web

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/auth"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
)

// AuthMiddleware extracts session ID from header and sets user context
func AuthMiddleware(authService auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		sessionID := c.GetHeader("X-Session-ID")
		if sessionID == "" {
			// No session ID provided, continue as anonymous user
			c.Set("user", &common.User{UserID: 0, Role: common.RoleUser})
			c.Next()
			return
		}

		user, err := authService.GetUserFromSession(c, sessionID)
		if err != nil {
			HandleError(c, err)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

// GetCommonContext creates a CommonContext from the gin context
func GetCommonContext(c *gin.Context) *common.CommonContext {
	user, exists := c.Get("user")
	if !exists {
		return common.NewCommonContext(c, &common.User{UserID: 0, Role: common.RoleUser})
	}

	userObj, ok := user.(*common.User)
	if !ok {
		return common.NewCommonContext(c, &common.User{UserID: 0, Role: common.RoleUser})
	}

	return common.NewCommonContext(c, userObj)
}

// RequireAuth middleware ensures user is authenticated
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			HandleError(c, apperror.ErrWrongAuth)
			c.Abort()
			return
		}

		userObj, ok := user.(*common.User)
		if !ok || userObj.UserID == 0 {
			HandleError(c, apperror.ErrWrongAuth)
			c.Abort()
			return
		}

		c.Next()
	}
}
