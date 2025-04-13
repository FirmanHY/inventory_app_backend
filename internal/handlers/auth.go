package handlers

import (
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

type AuthHandler struct {
	DB *gorm.DB
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	// Bind dan validasi request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
		return
	}

	// Cari user di database
	var user models.User
	result := h.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		utils.Error(c, http.StatusUnauthorized, constants.MsgInvalidCredentials, gin.H{
			"username": constants.MsgUsernameNotFound,
		})
		return
	}

	// Verifikasi password
	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
		utils.Error(c, http.StatusUnauthorized, constants.MsgInvalidCredentials, gin.H{
			"password": constants.MsgUsernameNotFound,
		})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.UserID, user.Role)
	if err != nil {
		utils.ServerError(c, constants.MsgTokenInvalid, err)
		return
	}

	resp := dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       user.UserID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	utils.Success(c, http.StatusOK, constants.MsgLoginSuccess, resp)
}

func (h *AuthHandler) CreateUser(c *gin.Context) {

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
		return
	}

	// Validasi role
	if !utils.IsValidRole(req.Role) {
		utils.BadRequest(c, constants.MsgInvalidRole, gin.H{
			"role": constants.MsgInvalidRoleValue,
		})
		return
	}

	// Cek username sudah ada
	var existingUser models.User
	if err := h.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.Error(c, http.StatusConflict, constants.MsgUsernameExists, gin.H{
			"username": constants.MsgUsernameRegistered,
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, constants.MsgInvalidRole, err)
		return
	}

	// Buat user baru
	newUser := models.User{
		UserID:   uuid.New().String(),
		Username: req.Username,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Role:     req.Role,
	}

	if err := h.DB.Create(&newUser).Error; err != nil {
		utils.ServerError(c, constants.MsgUserSaveFail, err)
		return
	}

	// Response
	resp := dto.UserResponse{
		ID:       newUser.UserID,
		Username: newUser.Username,
		FullName: newUser.FullName,
		Role:     newUser.Role,
	}

	utils.Success(c, http.StatusCreated, constants.MsgUserCreatedSuccess, resp)
}

func (h *AuthHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
		return
	}

	// Cari user yang akan diupdate
	var user models.User
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		utils.NotFound(c, constants.MsgUserNotFound)
		return
	}

	// Update username jika ada
	if req.Username != nil {
		// Cek username unik
		var existingUser models.User
		if err := h.DB.Where("username = ? AND user_id != ?", *req.Username, userID).First(&existingUser).Error; err == nil {
			utils.Error(c, http.StatusConflict, constants.MsgUsernameExists, gin.H{
				"username": constants.MsgUsernameRegistered,
			})
			return
		}
		user.Username = *req.Username
	}

	// Update full name jika ada
	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	// Update role jika ada
	if req.Role != nil {
		if !utils.IsValidRole(*req.Role) {
			utils.BadRequest(c, constants.MsgInvalidRole, gin.H{
				"role": constants.MsgInvalidRoleValue,
			})
			return
		}
		user.Role = *req.Role
	}

	// Simpan perubahan
	if err := h.DB.Save(&user).Error; err != nil {
		utils.ServerError(c, constants.MsgUserSaveFail, err)
		return
	}

	// Response
	resp := dto.UserResponse{
		ID:       user.UserID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}

	utils.Success(c, http.StatusOK, constants.MsgUserUpdateSuccess, resp)
}

func (h *AuthHandler) GetUsers(c *gin.Context) {
	// Parse query params
	page := utils.ParseInt(c.Query("page"), 1)
	limit := utils.ParseInt(c.Query("limit"), 10)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	var users []models.User

	// Query database
	if err := h.DB.Model(&models.User{}).Count(&total).Error; err != nil {
		utils.ServerError(c, constants.MsgInternalServerError, err)
		return
	}

	if err := h.DB.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		utils.ServerError(c, constants.MsgInternalServerError, err)
		return
	}

	// Map to response
	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:       user.UserID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		})
	}

	// Calculate pagination
	totalPages := total / int64(limit)
	if total%int64(limit) > 0 {
		totalPages++
	}

	resp := dto.UserListResponse{
		Data: userResponses,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			TotalData:  int(total),
			TotalPages: int(totalPages),
		},
	}

	utils.Success(c, http.StatusOK, constants.MsgUsersFetchSuccess, resp)
}

func (h *AuthHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := h.DB.Where("user_id = ?", userID).First(&user).Error; err != nil {
		utils.NotFound(c, constants.MsgUserNotFound)
		return
	}

	resp := dto.UserResponse{
		ID:       user.UserID,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
	}

	utils.Success(c, http.StatusOK, constants.MsgUserFetchSuccess, resp)
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	var req dto.ItemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, constants.MsgValidationFailed, err)
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

	query := h.DB.Model(&models.Item{}).
		Preload("Type").
		Preload("Unit").
		Order("created_at DESC")

	if req.Search != "" {
		query = query.Where("item_name LIKE ?", "%"+req.Search+"%")
	}

	if req.LowStockOnly {
		query = query.Where("stock < minimum_stock")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.ServerError(c, constants.MsgInternalServerError, err)
		return
	}

	var items []models.Item
	if err := query.Offset(offset).Limit(req.Limit).Find(&items).Error; err != nil {
		utils.ServerError(c, constants.MsgInternalServerError, err)
		return
	}

	var itemResponses []dto.ItemResponse
	for _, item := range items {
		itemResponses = append(itemResponses, dto.ItemResponse{
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
		})
	}

	totalPages := total / int64(req.Limit)
	if total%int64(req.Limit) > 0 {
		totalPages++
	}

	resp := dto.ItemListResponse{
		Data: itemResponses,
		Pagination: dto.Pagination{
			Page:       req.Page,
			Limit:      req.Limit,
			TotalData:  int(total),
			TotalPages: int(totalPages),
		},
	}

	utils.Success(c, http.StatusOK, constants.MsgItemsFetchSuccess, resp)
}
