package dto

type UnitResponse struct {
	UnitID   string `json:"unit_id"`
	UnitName string `json:"unit_name"`
}

type CreateUnitRequest struct {
	UnitName string `json:"unit_name" binding:"required,min=2"`
}

type UpdateUnitRequest struct {
	UnitName string `json:"unit_name" binding:"required,min=2"`
}

type UnitListRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
}

type UnitListResponse struct {
	Data       []UnitResponse `json:"data"`
	Pagination Pagination     `json:"pagination"`
}
