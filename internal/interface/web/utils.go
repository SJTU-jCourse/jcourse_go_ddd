package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"jcourse_go/internal/interface/dto"
	"jcourse_go/pkg/apperror"
)

func HandleError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*apperror.AppError); ok {
		response := dto.BaseResponse{
			Code: appErr.Code,
			Msg:  appErr.Message,
			Data: nil,
		}
		ctx.JSON(appErr.HTTPStatus(), response)
		return
	}

	// Handle generic errors
	response := dto.BaseResponse{
		Code: apperror.ErrDB.Code,
		Msg:  apperror.ErrDB.Message,
		Data: nil,
	}
	ctx.JSON(apperror.ErrDB.HTTPStatus(), response)
}

func HandleSuccess(ctx *gin.Context, data any) {
	response := dto.BaseResponse{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	ctx.JSON(http.StatusOK, response)
}

func HandleSuccessWithStatus(ctx *gin.Context, statusCode int, data any) {
	response := dto.BaseResponse{
		Code: 0,
		Msg:  "success",
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
