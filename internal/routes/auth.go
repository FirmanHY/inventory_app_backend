package routes

import (
	"inventory_app_backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(router *gin.Engine, h *handlers.AuthHandler) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/login", h.Login)
	}
}
