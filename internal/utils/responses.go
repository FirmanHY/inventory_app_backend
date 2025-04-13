package utils

import (
	constants "inventory_app_backend/internal/constant"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  string      `json:"status"`  // "success" atau "error"
	Message string      `json:"message"` // Pesan untuk developer/user
	Data    interface{} `json:"data"`    // Data payload
	Error   interface{} `json:"error"`   // Detail error (jika ada)
}

type ErrorResponse struct {
	Code    int         `json:"code"`    // Kode error internal
	Message string      `json:"message"` // Deskripsi error
	Details interface{} `json:"details"` // Detail error validasi
}

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Message: message,
		Data:    data,
		Error:   nil,
	})
}

func Error(c *gin.Context, statusCode int, message string, err interface{}) {
	var errorResponse ErrorResponse
	var validationErrors gin.H

	// Handle validation errors
	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		errorResponse.Code = statusCode
		errorResponse.Message = constants.MsgValidationError
		validationErrors = ParseValidationError(validationErrs)
		errorResponse.Details = validationErrors["validation"]
	} else if errMsg, ok := err.(error); ok {
		errorResponse = ErrorResponse{
			Code:    statusCode,
			Message: errMsg.Error(),
			Details: nil,
		}
	} else if errMsg, ok := err.(string); ok {
		errorResponse = ErrorResponse{
			Code:    statusCode,
			Message: errMsg,
			Details: nil,
		}
	} else {
		errorResponse = ErrorResponse{
			Code:    statusCode,
			Message: message,
			Details: err,
		}
	}

	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
		Data:    nil,
		Error:   errorResponse,
	})
}

// Response helpers untuk status code umum
func BadRequest(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusBadRequest, message, err)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, message, constants.MsgUnauthorizedError)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, message, constants.MsgForbiddenError)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message, constants.MsgNotFoundError)
}

func ServerError(c *gin.Context, message string, err interface{}) {
	Error(c, http.StatusInternalServerError, message, err)
}
