package handlers

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/dto"
	"inventory_app_backend/internal/models"
	"inventory_app_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SummaryHandler struct {
	DB *gorm.DB
}

func (h *SummaryHandler) GetInventorySummary(c *gin.Context) {
	summaryType := c.Query("type")

	var totalStock, totalIn, totalOut int64

	// Hitung total data barang dari semua item
	if err := h.DB.Model(&models.Item{}).Select("COUNT(*)").Scan(&totalStock).Error; err != nil {
		utils.ServerError(c, constants.MsgSummaryTotalFailed, err)
		return
	}

	// Hitung total barang masuk
	if err := h.DB.Model(&models.Transaction{}).
		Where("transaction_type = ?", constants.TransactionTypeIn).
		Select("SUM(quantity)").Scan(&totalIn).Error; err != nil {
		utils.ServerError(c, constants.MsgSummaryInFailed, err)
		return
	}

	// Hitung total barang keluar
	if err := h.DB.Model(&models.Transaction{}).
		Where("transaction_type = ?", constants.TransactionTypeOut).
		Select("SUM(quantity)").Scan(&totalOut).Error; err != nil {
		utils.ServerError(c, constants.MsgSummaryOutFailed, err)
		return
	}

	resp := dto.SummaryResponse{}

	switch summaryType {
	case constants.TransactionTypeIn:
		resp = dto.SummaryResponse{ItemsIn: totalIn}
	case constants.TransactionTypeOut:
		resp = dto.SummaryResponse{ItemsOut: totalOut}
	default:
		resp = dto.SummaryResponse{
			TotalItems: totalStock,
			ItemsIn:    totalIn,
			ItemsOut:   totalOut,
		}
	}

	utils.Success(c, http.StatusOK, constants.MsgSummaryFetchedSuccess, resp)
}
