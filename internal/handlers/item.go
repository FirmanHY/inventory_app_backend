package handlers

import (
	"fmt"
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/dto"
	"inventory_app_backend/internal/models"
	"inventory_app_backend/internal/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemHandler struct {
	DB *gorm.DB
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req dto.CreateItemRequest

	// Bind form data
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
		return
	}

	// Validate item type
	var itemType models.ItemType
	if err := h.DB.Where("type_id = ?", req.TypeID).First(&itemType).Error; err != nil {
		utils.BadRequest(c, constants.MsgInvalidItemType, gin.H{
			"type_id": constants.MsgItemTypeNotFound,
		})
		return
	}

	// Validate unit
	var unit models.Unit
	if err := h.DB.Where("unit_id = ?", req.UnitID).First(&unit).Error; err != nil {
		utils.BadRequest(c, constants.MsgInvalidUnit, gin.H{
			"unit_id": constants.MsgUnitNotFound,
		})
		return
	}

	// Upload image
	var imageURL string
	if req.Image != nil {
		url, validationErrors, err := utils.ValidateAndUploadImage(req.Image, "items")
		if err != nil {
			switch err.Error() {
			case "invalid_format":
				utils.BadRequest(c, constants.MsgInvalidImageFormat, validationErrors)
			case "too_large":
				utils.BadRequest(c, constants.MsgImageTooLarge, validationErrors)
			default:
				utils.ServerError(c, constants.MsgImageUploadFailed, err)
			}
			return
		}
		imageURL = url
	}

	// Create item
	newItem := models.Item{
		ItemID:       uuid.New().String(),
		TypeID:       req.TypeID,
		UnitID:       req.UnitID,
		ItemName:     req.ItemName,
		Stock:        0, // Stock awal selalu 0
		MinimumStock: req.MinimumStock,
		Image:        imageURL,
	}

	if err := h.DB.Create(&newItem).Error; err != nil {
		utils.ServerError(c, constants.MsgItemCreatedFailed, err)
		return
	}

	// Map ke response
	resp := dto.ItemDetailResponse{
		ItemResponse: dto.ItemResponse{
			ItemID:       newItem.ItemID,
			ItemName:     newItem.ItemName,
			TypeID:       newItem.TypeID,
			TypeName:     itemType.TypeName,
			UnitID:       newItem.UnitID,
			UnitName:     unit.UnitName,
			Stock:        newItem.Stock,
			MinimumStock: newItem.MinimumStock,
			Image:        newItem.Image,
			CreatedAt:    newItem.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    newItem.UpdatedAt.Format(time.RFC3339),
		},
		Type: dto.ItemTypeResponse{
			TypeID:   itemType.TypeID,
			TypeName: itemType.TypeName,
		},
		Unit: dto.UnitResponse{
			UnitID:   unit.UnitID,
			UnitName: unit.UnitName,
		},
	}

	utils.Success(c, http.StatusCreated, constants.MsgItemCreatedSuccess, resp)
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	itemID := c.Param("id")

	var req dto.UpdateItemRequest
	if err := c.ShouldBind(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
		return
	}

	// Cari item yang akan diupdate
	var item models.Item
	if err := h.DB.Where("item_id = ?", itemID).First(&item).Error; err != nil {
		utils.NotFound(c, constants.MsgItemTypeNotFound)
		return
	}

	// Update type jika ada
	if req.TypeID != nil {
		var itemType models.ItemType
		if err := h.DB.Where("type_id = ?", *req.TypeID).First(&itemType).Error; err != nil {
			utils.BadRequest(c, constants.MsgInvalidItemType, gin.H{
				"type_id": constants.MsgItemTypeNotFound,
			})
			return
		}
		item.TypeID = *req.TypeID
	}

	// Update unit jika ada
	if req.UnitID != nil {
		var unit models.Unit
		if err := h.DB.Where("unit_id = ?", *req.UnitID).First(&unit).Error; err != nil {
			utils.BadRequest(c, constants.MsgInvalidUnit, gin.H{
				"unit_id": constants.MsgUnitNotFound,
			})
			return
		}
		item.UnitID = *req.UnitID
	}

	// Update nama item jika ada
	if req.ItemName != nil {
		item.ItemName = *req.ItemName
	}

	// Update minimum stock jika ada
	if req.MinimumStock != nil {
		item.MinimumStock = *req.MinimumStock
	}

	// Update gambar jika ada
	if req.Image != nil {
		url, validationErrors, err := utils.ValidateAndUploadImage(req.Image, "items")
		if err != nil {
			switch err.Error() {
			case "invalid_format":
				utils.BadRequest(c, constants.MsgInvalidImageFormat, validationErrors)
			case "too_large":
				utils.BadRequest(c, constants.MsgImageTooLarge, validationErrors)
			default:
				utils.ServerError(c, constants.MsgImageUploadFailed, err)
			}
			return
		}

		item.Image = url
	}

	// Simpan perubahan
	if err := h.DB.Save(&item).Error; err != nil {
		utils.ServerError(c, constants.MsgItemUpdatedFailed, err)
		return
	}

	// Get updated data dengan relasi
	var updatedItem models.Item
	if err := h.DB.Preload("Type").Preload("Unit").
		First(&updatedItem, "item_id = ?", itemID).Error; err != nil {
		utils.ServerError(c, constants.MsgGetUpdatedItemFailed, err)
		return
	}

	// Map ke response
	resp := dto.ItemDetailResponse{
		ItemResponse: dto.ItemResponse{
			ItemID:       updatedItem.ItemID,
			ItemName:     updatedItem.ItemName,
			TypeID:       updatedItem.TypeID,
			TypeName:     updatedItem.Type.TypeName,
			UnitID:       updatedItem.UnitID,
			UnitName:     updatedItem.Unit.UnitName,
			Stock:        updatedItem.Stock,
			MinimumStock: updatedItem.MinimumStock,
			Image:        updatedItem.Image,
			CreatedAt:    updatedItem.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    updatedItem.UpdatedAt.Format(time.RFC3339),
		},
		Type: dto.ItemTypeResponse{
			TypeID:   updatedItem.Type.TypeID,
			TypeName: updatedItem.Type.TypeName,
		},
		Unit: dto.UnitResponse{
			UnitID:   updatedItem.Unit.UnitID,
			UnitName: updatedItem.Unit.UnitName,
		},
	}

	utils.Success(c, http.StatusOK, constants.MsgItemUpdatedSuccess, resp)
}

func (h *ItemHandler) GetItemByID(c *gin.Context) {
	itemID := c.Param("id")

	var item models.Item
	if err := h.DB.Preload("Type").Preload("Unit").
		First(&item, "item_id = ?", itemID).Error; err != nil {
		utils.NotFound(c, constants.MsgItemNotFound)
		return
	}

	resp := dto.ItemDetailResponse{
		ItemResponse: dto.ItemResponse{
			ItemID:       item.ItemID,
			ItemName:     item.ItemName,
			TypeID:       item.TypeID,
			TypeName:     item.Type.TypeName,
			UnitID:       item.UnitID,
			UnitName:     item.Unit.UnitName,
			Stock:        item.Stock,
			MinimumStock: item.MinimumStock,
			Image:        item.Image,
			CreatedAt:    item.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    item.UpdatedAt.Format(time.RFC3339),
		},
		Type: dto.ItemTypeResponse{
			TypeID:   item.Type.TypeID,
			TypeName: item.Type.TypeName,
		},
		Unit: dto.UnitResponse{
			UnitID:   item.Unit.UnitID,
			UnitName: item.Unit.UnitName,
		},
	}

	utils.Success(c, http.StatusOK, constants.MsgItemFetchSuccess, resp)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	itemID := c.Param("id")

	// Cek apakah item ada
	var item models.Item
	if err := h.DB.Where("item_id = ?", itemID).First(&item).Error; err != nil {
		utils.NotFound(c, constants.MsgItemNotFound)
		return
	}

	// Cek apakah item digunakan di transaksi
	var transactionCount int64
	if err := h.DB.Model(&models.Transaction{}).Where("item_id = ?", itemID).Count(&transactionCount).Error; err != nil {
		utils.ServerError(c, constants.MsgItemDeleteFailed, err)
		return
	}

	if transactionCount > 0 {
		utils.Error(c, http.StatusConflict, constants.MsgItemInUse, gin.H{
			"item_id": fmt.Sprintf(constants.MsgItemInUseDetail, transactionCount),
		})
		return
	}

	// Hapus item
	if err := h.DB.Delete(&item).Error; err != nil {
		utils.ServerError(c, constants.MsgItemDeleteFailed, err)
		return
	}

	utils.Success(c, http.StatusOK, constants.MsgItemDeletedSuccess, nil)
}

func (h *ItemHandler) GetLowStockItems(c *gin.Context) {
	var items []models.Item
	if err := h.DB.Preload("Type").
		Where("stock < minimum_stock").
		Find(&items).Error; err != nil {
		utils.ServerError(c, constants.MsgFailedFetchItems, err)
		return
	}

	var result []dto.ItemResponse
	for _, item := range items {
		result = append(result, dto.ItemResponse{
			ItemID:   item.ItemID,
			ItemName: item.ItemName,
			TypeID:   item.TypeID,
			TypeName: item.Type.TypeName,
			Stock:    item.Stock,
		})
	}

	utils.Success(c, http.StatusOK, constants.MsgItemsFetchSuccess, result)
}
