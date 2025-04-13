package middleware

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			switch err.Type {
			case gin.ErrorTypeBind:
				utils.BadRequest(c, constants.MsgInvalidRequest, err.Err)
			case gin.ErrorTypeRender:
				utils.ServerError(c, constants.MsgRenderFailed, err.Err)
			default:
				utils.ServerError(c, constants.MsgInternalServerError, err.Err)
			}

			c.Abort()
			return
		}
	}
}

// Recovery middleware untuk handle panic
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			utils.ServerError(c, constants.MsgInternalServerError, err)
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	})
}
