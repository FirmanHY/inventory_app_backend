package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupItemRoutes(router *gin.Engine, h *handlers.ItemHandler) {

	itemRoutes := router.Group("/items")
	itemRoutes.Use(middleware.Auth(), middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		itemRoutes.GET("", h.GetAllItems)
		itemRoutes.GET("/:id", h.GetItemByID)
		itemRoutes.POST("", h.CreateItem)
		itemRoutes.PUT("/:id", h.UpdateItem)
		itemRoutes.DELETE("/:id", h.DeleteItem)
	}
}
