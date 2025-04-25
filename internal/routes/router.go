package routes

import (
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *handlers.AuthHandler,
	itemHandler *handlers.ItemHandler,
	itemTypeHandler *handlers.ItemTypeHandler,
	unitHandler *handlers.UnitHandler,
	transactionHandler *handlers.TransactionHandler,
	reportHandler *handlers.ReportHandler,
	summaryHandler *handlers.SummaryHandler,

) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(
		middleware.Recovery(),
		gin.Logger(),
		middleware.ErrorHandler(),
	)

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Setup route groups
	setupAuthRoutes(router, authHandler)
	setupAdminRoutes(router, authHandler)
	setupItemRoutes(router, itemHandler)
	setupMasterDataRoutes(router, itemTypeHandler, unitHandler)
	setupTransactionRoutes(router, transactionHandler)
	setupReportRoutes(router, reportHandler)
	setupSummaryRoutes(router, summaryHandler)
	return router
}
