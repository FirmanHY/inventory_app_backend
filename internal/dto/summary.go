package dto

type SummaryResponse struct {
	TotalBarang  int64 `json:"total_barang,omitempty"`
	BarangMasuk  int64 `json:"barang_masuk,omitempty"`
	BarangKeluar int64 `json:"barang_keluar,omitempty"`
}
