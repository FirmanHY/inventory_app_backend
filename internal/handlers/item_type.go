package handlers

import (
	"fmt"
	constant "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/dto"
	"inventory_app_backend/internal/models"
	"inventory_app_backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemTypeHandler struct {
	DB *gorm.DB
}

func (h *ItemTypeHandler) CreateItemType(c *gin.Context) {
	var req dto.CreateItemTypeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Cek duplikat type name
	var existingType models.ItemType
	if err := h.DB.Where("type_name = ?", req.TypeName).First(&existingType).Error; err == nil {
		utils.Error(c, http.StatusConflict, constant.MsgItemTypeExists, gin.H{
			"type_name": constant.MsgItemTypeExistsDetail,
		})
		return
	}

	// Buat item type baru
	newItemType := models.ItemType{
		TypeID:    uuid.New().String(),
		TypeName:  req.TypeName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.DB.Create(&newItemType).Error; err != nil {
		utils.ServerError(c, constant.MsgItemTypeCreateFailed, err)
		return
	}

	resp := dto.ItemTypeResponse{
		TypeID:   newItemType.TypeID,
		TypeName: newItemType.TypeName,
	}

	utils.Success(c, http.StatusCreated, constant.MsgItemTypeCreatedSuccess, resp)
}

func (h *ItemTypeHandler) UpdateItemType(c *gin.Context) {
	typeID := c.Param("id")

	var req dto.UpdateItemTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Cek apakah type ada
	var itemType models.ItemType
	if err := h.DB.Where("type_id = ?", typeID).First(&itemType).Error; err != nil {
		utils.NotFound(c, constant.MsgItemTypeNotFound)
		return
	}

	// Cek duplikat nama type
	var existingType models.ItemType
	if err := h.DB.Where("type_name = ? AND type_id != ?", req.TypeName, typeID).First(&existingType).Error; err == nil {
		utils.Error(c, http.StatusConflict, constant.MsgItemTypeExists, gin.H{
			"type_name": constant.MsgItemTypeExistsDetail,
		})
		return
	}

	// Update data
	itemType.TypeName = req.TypeName
	itemType.UpdatedAt = time.Now()

	if err := h.DB.Save(&itemType).Error; err != nil {
		utils.ServerError(c, constant.MsgItemTypeUpdateFailed, err)
		return
	}

	resp := dto.ItemTypeResponse{
		TypeID:   itemType.TypeID,
		TypeName: itemType.TypeName,
	}

	utils.Success(c, http.StatusOK, constant.MsgItemTypeUpdatedSuccess, resp)
}

func (h *ItemTypeHandler) GetAllItemTypes(c *gin.Context) {
	var req dto.ItemTypeListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Set default values
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	// Build query
	query := h.DB.Model(&models.ItemType{}).Order("created_at DESC")

	// Add search filter
	if req.Search != "" {
		query = query.Where("type_name LIKE ?", "%"+req.Search+"%")
	}

	// Get total data
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	// Get paginated data
	var itemTypes []models.ItemType
	if err := query.Offset(offset).Limit(req.Limit).Find(&itemTypes).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	// Map to response
	var typeResponses []dto.ItemTypeResponse
	for _, itemType := range itemTypes {
		typeResponses = append(typeResponses, dto.ItemTypeResponse{
			TypeID:   itemType.TypeID,
			TypeName: itemType.TypeName,
		})
	}

	// Calculate pagination
	totalPages := total / int64(req.Limit)
	if total%int64(req.Limit) > 0 {
		totalPages++
	}

	resp := dto.ItemTypeListResponse{
		Data: typeResponses,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalData:  int(total),
			TotalPages: int(totalPages),
		},
	}

	utils.Success(c, http.StatusOK, constant.MsgItemTypesFetchSuccess, resp)
}

func (h *ItemTypeHandler) DeleteItemType(c *gin.Context) {
	typeID := c.Param("id")

	// Cek apakah type ada
	var itemType models.ItemType
	if err := h.DB.Where("type_id = ?", typeID).First(&itemType).Error; err != nil {
		utils.NotFound(c, constant.MsgItemTypeNotFound)
		return
	}

	// Cek apakah type digunakan di items
	var count int64
	if err := h.DB.Model(&models.Item{}).Where("type_id = ?", typeID).Count(&count).Error; err != nil {
		utils.ServerError(c, constant.MsgItemTypeDeleteFailed, err)
		return
	}

	if count > 0 {
		utils.Error(c, http.StatusConflict, constant.MsgItemTypeInUse, gin.H{
			"type_id": fmt.Sprintf(constant.MsgItemTypeInUseDetail, count),
		})
		return
	}

	// Hapus type
	if err := h.DB.Delete(&itemType).Error; err != nil {
		utils.ServerError(c, constant.MsgItemTypeDeleteFailed, err)
		return
	}

	utils.Success(c, http.StatusOK, constant.MsgItemTypeDeletedSuccess, nil)
}
