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
