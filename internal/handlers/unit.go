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

type UnitHandler struct {
	DB *gorm.DB
}

func (h *UnitHandler) CreateUnit(c *gin.Context) {
	var req dto.CreateUnitRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Cek duplikat
	var existingUnit models.Unit
	if err := h.DB.Where("unit_name = ?", req.UnitName).First(&existingUnit).Error; err == nil {
		utils.Error(c, http.StatusConflict, constant.MsgUnitExists, gin.H{
			"unit_name": constant.MsgUnitExistsDetail,
		})
		return
	}

	newUnit := models.Unit{
		UnitID:    uuid.New().String(),
		UnitName:  req.UnitName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.DB.Create(&newUnit).Error; err != nil {
		utils.ServerError(c, constant.MsgUnitCreateFailed, err)
		return
	}

	resp := dto.UnitResponse{
		UnitID:   newUnit.UnitID,
		UnitName: newUnit.UnitName,
	}

	utils.Success(c, http.StatusCreated, constant.MsgUnitCreatedSuccess, resp)
}

func (h *UnitHandler) UpdateUnit(c *gin.Context) {
	unitID := c.Param("id")

	var req dto.UpdateUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Cek unit
	var unit models.Unit
	if err := h.DB.Where("unit_id = ?", unitID).First(&unit).Error; err != nil {
		utils.NotFound(c, constant.MsgUnitNotFound)
		return
	}

	// Cek duplikat
	var existingUnit models.Unit
	if err := h.DB.Where("unit_name = ? AND unit_id != ?", req.UnitName, unitID).First(&existingUnit).Error; err == nil {
		utils.Error(c, http.StatusConflict, constant.MsgUnitExists, gin.H{
			"unit_name": constant.MsgUnitExistsDetail,
		})
		return
	}

	unit.UnitName = req.UnitName
	unit.UpdatedAt = time.Now()

	if err := h.DB.Save(&unit).Error; err != nil {
		utils.ServerError(c, constant.MsgUnitUpdateFailed, err)
		return
	}

	resp := dto.UnitResponse{
		UnitID:   unit.UnitID,
		UnitName: unit.UnitName,
	}

	utils.Success(c, http.StatusOK, constant.MsgUnitUpdatedSuccess, resp)
}

func (h *UnitHandler) GetAllUnits(c *gin.Context) {
	var req dto.UnitListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, constant.MsgValidationFailed, err)
		return
	}

	// Set default
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 10
	}
	offset := (req.Page - 1) * req.Limit

	// Query
	query := h.DB.Model(&models.Unit{}).Order("created_at DESC")

	if req.Search != "" {
		query = query.Where("unit_name LIKE ?", "%"+req.Search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	var units []models.Unit
	if err := query.Offset(offset).Limit(req.Limit).Find(&units).Error; err != nil {
		utils.ServerError(c, constant.MsgInternalServerError, err)
		return
	}

	// Mapping
	var unitResponses []dto.UnitResponse
	for _, unit := range units {
		unitResponses = append(unitResponses, dto.UnitResponse{
			UnitID:   unit.UnitID,
			UnitName: unit.UnitName,
		})
	}

	// Pagination
	totalPages := total / int64(req.Limit)
	if total%int64(req.Limit) > 0 {
		totalPages++
	}

	resp := dto.UnitListResponse{
		Data: unitResponses,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalData:  int(total),
			TotalPages: int(totalPages),
		},
	}

	utils.Success(c, http.StatusOK, constant.MsgUnitsFetchSuccess, resp)
}

func (h *UnitHandler) DeleteUnit(c *gin.Context) {
	unitID := c.Param("id")

	var unit models.Unit
	if err := h.DB.Where("unit_id = ?", unitID).First(&unit).Error; err != nil {
		utils.NotFound(c, constant.MsgUnitNotFound)
		return
	}

	// Cek apakah unit digunakan di items
	var count int64
	if err := h.DB.Model(&models.Item{}).Where("unit_id = ?", unitID).Count(&count).Error; err != nil {
		utils.ServerError(c, constant.MsgItemTypeDeleteFailed, err)
		return
	}

	if count > 0 {
		utils.Error(c, http.StatusConflict, constant.MsgItemTypeInUse, gin.H{
			"unit_id": fmt.Sprintf(constant.MsgUnitInUseDetail, count),
		})
		return
	}

	if err := h.DB.Delete(&unit).Error; err != nil {
		utils.ServerError(c, constant.MsgUnitDeleteFailed, err)
		return
	}

	utils.Success(c, http.StatusOK, constant.MsgUnitDeletedSuccess, nil)
}
