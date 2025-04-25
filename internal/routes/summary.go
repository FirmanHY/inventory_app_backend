package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupSummaryRoutes(router *gin.Engine, h *handlers.SummaryHandler) {
	summary := router.Group("/summary")
	summary.Use(middleware.Auth(), middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin, constants.RoleWarehouseManager))
	{
		summary.GET("/inventory", h.GetInventorySummary)
	}
}
