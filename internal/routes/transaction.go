package routes

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func setupTransactionRoutes(router *gin.Engine, h *handlers.TransactionHandler) {
	transactionRoutes := router.Group("/transactions")
	transactionRoutes.Use(middleware.Auth())

	readRoutes := transactionRoutes.Group("")
	readRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin, constants.RoleWarehouseManager))
	{
		readRoutes.GET("", h.GetAllTransactions)
	}

	writeRoutes := transactionRoutes.Group("")
	writeRoutes.Use(middleware.RoleAllowed(constants.RoleAdmin, constants.RoleWarehouseAdmin))
	{
		writeRoutes.POST("", h.CreateTransaction)
		writeRoutes.DELETE("/:id", h.DeleteTransaction)
	}
}
