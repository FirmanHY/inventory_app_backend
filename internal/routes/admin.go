package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupAdminRoutes(router *gin.Engine, h *handlers.AuthHandler) {
	adminGroup := router.Group("/admin")
	adminGroup.Use(
		middleware.Auth(),
		middleware.RoleAllowed(constants.RoleAdmin),
	)
	{
		adminGroup.GET("/users", h.GetUsers)
		adminGroup.GET("/users/:id", h.GetUserByID)
		adminGroup.POST("/users", h.CreateUser)
		adminGroup.PUT("/users/:id", h.UpdateUser)
	}
}
