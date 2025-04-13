package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupTransactionRoutes(router *gin.Engine, h *handlers.TransactionHandler) {

	transactionRoutes := router.Group("/transactions")
	transactionRoutes.Use(middleware.Auth(), middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		transactionRoutes.GET("", h.GetAllTransactions)
		transactionRoutes.POST("", h.CreateTransaction)
		transactionRoutes.DELETE("/:id", h.DeleteTransaction)
	}
}
