package middleware

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/utils"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			utils.Unauthorized(c, constants.MsgAuthHeaderRequired)
			c.Abort()
			return
		}

		// Validate JWT
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.Unauthorized(c, constants.MsgTokenInvalid)
			c.Abort()
			return
		}

		// Set user context
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
