package dto

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalData  int `json:"total_data"`
	TotalPages int `json:"total_pages"`
}
