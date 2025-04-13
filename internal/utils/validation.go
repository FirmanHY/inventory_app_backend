package utils

import (
	constants "inventory_app_backend/internal/constant"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// ParseValidationError mengembalikan map dari error validasi
func ParseValidationError(err error) gin.H {
	errors := gin.H{}
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := strings.ToLower(e.Field())
			switch e.Tag() {
			case "required":
				errors[field] = constants.MsgFieldRequired
			case "min":
				if field == "username" {
					errors[field] = constants.MsgUsernameMinLength
				} else if field == "password" {
					errors[field] = constants.MsgPasswordMinLength
				}
			default:
				errors[field] = constants.MsgInvalidFieldFormat
			}
		}
	}
	return gin.H{"validation": errors}
}

// Validasi role users
func IsValidRole(role string) bool {
	switch role {
	case constants.RoleAdmin,
		constants.RoleWarehouseAdmin,
		constants.RoleWarehouseManager:
		return true
	default:
		return false
	}
}
