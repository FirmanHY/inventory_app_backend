package dto

import "mime/multipart"

type CreateItemRequest struct {
	ItemName     string                `form:"item_name" binding:"required"`
	TypeID       string                `form:"type_id" binding:"required,uuid"`
	UnitID       string                `form:"unit_id" binding:"required,uuid"`
	MinimumStock int                   `form:"minimum_stock" binding:"min=0"`
	Image        *multipart.FileHeader `form:"image"`
}

type UpdateItemRequest struct {
	ItemName     *string               `form:"item_name" binding:"omitempty"`
	TypeID       *string               `form:"type_id" binding:"omitempty,uuid"`
	UnitID       *string               `form:"unit_id" binding:"omitempty,uuid"`
	MinimumStock *int                  `form:"minimum_stock" binding:"omitempty,min=0"`
	Image        *multipart.FileHeader `form:"image"`
}

type ItemListRequest struct {
	Page         int    `form:"page" binding:"omitempty,min=1"`
	Limit        int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search       string `form:"search"`
	LowStockOnly bool   `form:"low_stock_only"`
}

type ItemResponse struct {
	ItemID       string `json:"item_id"`
	ItemName     string `json:"item_name"`
	TypeID       string `json:"type_id"`
	TypeName     string `json:"type_name"`
	UnitID       string `json:"unit_id"`
	UnitName     string `json:"unit_name"`
	Stock        int    `json:"stock"`
	MinimumStock int    `json:"minimum_stock"`
	Image        string `json:"image"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ItemDetailResponse struct {
	ItemResponse
	Type ItemTypeResponse `json:"type"`
	Unit UnitResponse     `json:"unit"`
}

type ItemListResponse struct {
	Data       []ItemResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
