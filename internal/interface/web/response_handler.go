package web

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/pkg/apperror"
)

const (
	SuccessCode     = 0
	ResponseSuccess = "success"
)

func HandleError(ctx *gin.Context, err error) {
	// Use the new error handler for consistent responses
	errorResponse := apperror.HandleError(ctx.Request.Context(), err)

	response := dto.BaseResponse{
		Code: errorResponse.Code,
		Msg:  errorResponse.Message,
		Data: nil,
	}

	// Use user-friendly message if available
	if errorResponse.UserMessage != "" {
		response.Msg = errorResponse.UserMessage
	}

	// Determine HTTP status from error category
	var appErr *apperror.AppError
	if errors.As(err, &appErr) {
		ctx.JSON(appErr.HTTPStatus(), response)
	} else {
		ctx.JSON(http.StatusInternalServerError, response)
	}
}

func HandleSuccess(ctx *gin.Context, data any) {
	response := dto.BaseResponse{
		Code: SuccessCode,
		Msg:  ResponseSuccess,
		Data: data,
	}
	ctx.JSON(http.StatusOK, response)
}

func HandleSuccessWithStatus(ctx *gin.Context, statusCode int, data any) {
	response := dto.BaseResponse{
		Code: SuccessCode,
		Msg:  ResponseSuccess,
		Data: data,
	}
	ctx.JSON(statusCode, response)
}

func HandleValidationError(ctx *gin.Context, message string) {
	response := dto.BaseResponse{
		Code: apperror.ErrWrongInput.Code,
		Msg:  message,
		Data: nil,
	}
	ctx.JSON(apperror.ErrWrongInput.HTTPStatus(), response)
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
