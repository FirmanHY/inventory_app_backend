package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupItemRoutes(router *gin.Engine, h *handlers.ItemHandler) {
	itemRoutes := router.Group("/items")
	itemRoutes.Use(middleware.Auth())

	// Read-only routes accessible to admin, warehouse_admin, and warehouse_manager
	readRoutes := itemRoutes.Group("")
	readRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin, constants.RoleWarehouseManager))
	{
		readRoutes.GET("", h.GetAllItems)
		readRoutes.GET("/:id", h.GetItemByID)
		readRoutes.GET("/low-stock", h.GetLowStockItems)
	}

	// Write routes accessible only to admin and warehouse_admin
	writeRoutes := itemRoutes.Group("")
	writeRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		writeRoutes.POST("", h.CreateItem)
		writeRoutes.PUT("/:id", h.UpdateItem)
		writeRoutes.DELETE("/:id", h.DeleteItem)
	}
}
