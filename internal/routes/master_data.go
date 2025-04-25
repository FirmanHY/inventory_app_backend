package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupMasterDataRoutes(router *gin.Engine, h *handlers.ItemTypeHandler, u *handlers.UnitHandler) {
	masterDataRoutes := router.Group("/master-data")
	masterDataRoutes.Use(middleware.Auth())

	// Read-only routes accessible to admin, warehouse_admin, and warehouse_manager
	readRoutes := masterDataRoutes.Group("")
	readRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin, constants.RoleWarehouseManager))
	{
		readRoutes.GET("/item-types", h.GetAllItemTypes)
		readRoutes.GET("/units", u.GetAllUnits)
	}

	// Write routes accessible only to admin and warehouse_admin
	writeRoutes := masterDataRoutes.Group("")
	writeRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		writeRoutes.POST("/item-types", h.CreateItemType)
		writeRoutes.PUT("/item-types/:id", h.UpdateItemType)
		writeRoutes.DELETE("/item-types/:id", h.DeleteItemType)
		writeRoutes.POST("/units", u.CreateUnit)
		writeRoutes.PUT("/units/:id", u.UpdateUnit)
		writeRoutes.DELETE("/units/:id", u.DeleteUnit)
	}
}
