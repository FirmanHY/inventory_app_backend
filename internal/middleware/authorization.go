package middleware

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func RoleAllowed(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.Forbidden(c, constants.MsgForbiddenError)
			c.Abort()
			return
		}

		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		utils.Forbidden(c, constants.MsgForbiddenError)
		c.Abort()
	}
}
