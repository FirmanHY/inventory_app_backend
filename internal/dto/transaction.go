package dto

import "time"

type CreateTransactionRequest struct {
	ItemID          string    `json:"item_id" binding:"required,uuid"`
	Date            time.Time `json:"date" binding:"required" time_format:"2006-01-02"`
	Quantity        int       `json:"quantity" binding:"required,min=1"`
	TransactionType string    `json:"transaction_type" binding:"required,oneof=in out"`
	Description     string    `json:"description"`
}

type TransactionListRequest struct {
	Page       int       `form:"page" binding:"omitempty,min=1"`
	Limit      int       `form:"limit" binding:"omitempty,min=1,max=100"`
	Search     string    `form:"search"`
	StartDate  time.Time `form:"start_date" time_format:"2006-01-02"`
	EndDate    time.Time `form:"end_date" time_format:"2006-01-02"`
	TypeFilter string    `form:"type" binding:"omitempty,oneof=in out"`
}

type TransactionResponse struct {
	TransactionID   string    `json:"transaction_id"`
	ItemID          string    `json:"item_id"`
	ItemName        string    `json:"item_name"`
	Date            time.Time `json:"date"`
	Quantity        int       `json:"quantity"`
	TransactionType string    `json:"transaction_type"`
	Description     string    `json:"description"`
	CurrentStock    int       `json:"current_stock"`
	CreatedAt       time.Time `json:"created_at"`
}

type TransactionListResponse struct {
	Data       []TransactionResponse `json:"data"`
	Pagination Pagination            `json:"pagination"`
}

type DeleteTransactionResponse struct {
	TransactionID string `json:"transaction_id"`
	ItemID        string `json:"item_id"`
	CurrentStock  int    `json:"current_stock"`
	Warning       string `json:"warning,omitempty"`
}
