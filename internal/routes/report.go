package routes

import (
	"inventory_app_backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func setupReportRoutes(router *gin.Engine, h *handlers.ReportHandler) {

	reportRoutes := router.Group("/reports")

	{
		reportRoutes.GET("/items", h.GenerateItemReport)
		reportRoutes.GET("/transactions", h.GenerateTransactionReport)
	}
}
