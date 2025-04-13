package dto

type CreateItemTypeRequest struct {
	TypeName string `json:"type_name" binding:"required,min=3"`
}

type UpdateItemTypeRequest struct {
	TypeName string `json:"type_name" binding:"required,min=3"`
}

type ItemTypeResponse struct {
	TypeID   string `json:"type_id"`
	TypeName string `json:"type_name"`
}

type ItemTypeListRequest struct {
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
	Search string `form:"search"`
}

type ItemTypeListResponse struct {
	Data       []ItemTypeResponse `json:"data"`
	Pagination Pagination         `json:"pagination"`
}
