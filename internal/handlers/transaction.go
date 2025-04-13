package handlers

import (
	"fmt"
	constant "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/dto"
	"inventory_app_backend/internal/models"
	"inventory_app_backend/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	DB *gorm.DB
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req dto.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Cek item exists
	var item models.Item
	if err := h.DB.First(&item, "item_id = ?", req.ItemID).Error; err != nil {
		utils.NotFound(c, constant.MsgItemNotFound)
		return
	}

	// Untuk transaksi keluar, cek stok cukup
	if req.TransactionType == constant.TransactionTypeOut && item.Stock < req.Quantity {
		utils.Error(c, http.StatusBadRequest, constant.MsgInsufficientStock, gin.H{
			"current_stock": item.Stock,
			"required":      req.Quantity,
		})
		return
	}

	// Dapatkan user ID dari context
	userID, exists := c.Get("userID")
	if !exists {
		utils.Unauthorized(c, constant.MsgInvalidSession)
		return
	}

	// Buat transaksi
	newTransaction := models.Transaction{
		TransactionID:   uuid.New().String(),
		ItemID:          req.ItemID,
		Date:            req.Date,
		Quantity:        req.Quantity,
		TransactionType: req.TransactionType,
		Description:     req.Description,
		UserID:          userID.(string),
	}

	if err := h.DB.Create(&newTransaction).Error; err != nil {
		utils.ServerError(c, constant.MsgTransactionCreatedFailed, err)
		return
	}

	// Dapatkan stok terbaru setelah trigger dijalankan
	var updatedItem models.Item
	h.DB.First(&updatedItem, "item_id = ?", req.ItemID)

	resp := dto.TransactionResponse{
		TransactionID:   newTransaction.TransactionID,
		ItemID:          newTransaction.ItemID,
		ItemName:        item.ItemName,
		Date:            newTransaction.Date,
		Quantity:        newTransaction.Quantity,
		TransactionType: newTransaction.TransactionType,
		Description:     newTransaction.Description,
		CurrentStock:    updatedItem.Stock,
		CreatedAt:       newTransaction.CreatedAt,
	}

	utils.Success(c, http.StatusCreated, constant.MsgTransactionCreatedSuccess, resp)
}

func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	var req dto.TransactionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	fmt.Println("Parsed StartDate:", req.StartDate)
	fmt.Println("Parsed EndDate:", req.EndDate)
	fmt.Println("Is StartDate zero?", req.StartDate.IsZero())
	fmt.Println("Is EndDate zero?", req.EndDate.IsZero())

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	// Build query
	query := h.DB.Model(&models.Transaction{}).
		Preload("Item").
		Preload("User").
		Order("transactions.date DESC")

		// Apply filters
	if req.Search != "" {
		// Gunakan subquery untuk menghindari konflik join
		subQuery := h.DB.Model(&models.Item{}).
			Select("item_id").
			Where("item_name LIKE ?", "%"+req.Search+"%")

		query = query.Where("transactions.item_id IN (?)", subQuery)
	}

	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		if req.StartDate.After(req.EndDate) {

			req.StartDate, req.EndDate = req.EndDate, req.StartDate
		}

		query = query.Where("transactions.date BETWEEN ? AND ?",
			req.StartDate.Format("2006-01-02"),
			req.EndDate.Format("2006-01-02"))
	}

	if req.TypeFilter != "" {
		query = query.Where("transaction_type = ?", req.TypeFilter)
	}

	// Get total data
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	// Get paginated data
	var transactions []models.Transaction
	if err := query.Offset(offset).Limit(req.Limit).Find(&transactions).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	// Map to response
	var transactionResponses []dto.TransactionResponse
	for _, t := range transactions {
		transactionResponses = append(transactionResponses, dto.TransactionResponse{
			TransactionID:   t.TransactionID,
			ItemID:          t.ItemID,
			ItemName:        t.Item.ItemName,
			Date:            t.Date,
			Quantity:        t.Quantity,
			TransactionType: t.TransactionType,
			Description:     t.Description,
			CurrentStock:    t.Item.Stock,
			CreatedAt:       t.CreatedAt,
		})
	}

	// Calculate pagination
	totalPages := total / int64(req.Limit)
	if total%int64(req.Limit) > 0 {
		totalPages++
	}

	resp := dto.TransactionListResponse{
		Data: transactionResponses,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalData:  int(total),
			TotalPages: int(totalPages),
		},
	}

	utils.Success(c, http.StatusOK, constant.MsgTransactionFetchedSuccess, resp)
}

func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	transactionID := c.Param("id")

	// Cek transaksi exists
	var transaction models.Transaction
	if err := h.DB.Preload("Item").First(&transaction, "transaction_id = ?", transactionID).Error; err != nil {
		utils.NotFound(c, constant.MsgTransactionNotFound)
		return
	}

	// Simpan stok sebelum dihapus
	originalStock := transaction.Item.Stock

	// Hapus transaksi
	if err := h.DB.Delete(&transaction).Error; err != nil {
		utils.ServerError(c, constant.MsgTransactionDeleteFailed, err)
		return
	}

	// Get stok terbaru
	var updatedItem models.Item
	h.DB.First(&updatedItem, "item_id = ?", transaction.ItemID)

	response := dto.DeleteTransactionResponse{
		TransactionID: transaction.TransactionID,
		ItemID:        transaction.ItemID,
		CurrentStock:  updatedItem.Stock,
	}

	if transaction.TransactionType == "in" {
		expectedStock := originalStock - transaction.Quantity

		if expectedStock < 0 && updatedItem.Stock == 0 {
			response.Warning = constant.MsgTransactionAdjustStock
		}
	}

	utils.Success(c, http.StatusOK, constant.MsgTransactionDeletedSuccess, response)
}
