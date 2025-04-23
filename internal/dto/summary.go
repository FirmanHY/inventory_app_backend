package dto

type SummaryResponse struct {
	TotalItems int64 `json:"total_items,omitempty"`
	ItemsIn    int64 `json:"items_in,omitempty"`
	ItemsOut   int64 `json:"items_out,omitempty"`
}
