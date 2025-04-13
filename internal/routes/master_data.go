package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupMasterDataRoutes(router *gin.Engine, h *handlers.ItemTypeHandler, u *handlers.UnitHandler) {

	masterDataRoutes := router.Group("/master-data")
	masterDataRoutes.Use(middleware.Auth(), middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		masterDataRoutes.POST("/item-types", h.CreateItemType)
		masterDataRoutes.PUT("/item-types/:id", h.UpdateItemType)
		masterDataRoutes.GET("/item-types", h.GetAllItemTypes)
		masterDataRoutes.DELETE("/item-types/:id", h.DeleteItemType)
		// Unit routes
		masterDataRoutes.GET("/units", u.GetAllUnits)
		masterDataRoutes.POST("/units", u.CreateUnit)
		masterDataRoutes.PUT("/units/:id", u.UpdateUnit)
		masterDataRoutes.DELETE("/units/:id", u.DeleteUnit)
	}
}
