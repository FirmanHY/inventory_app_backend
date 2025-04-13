package constants

// ========================
// AUTHENTICATION MESSAGES
// ========================
const (
	MsgLoginSuccess        = "Login berhasil"
	MsgInvalidCredentials  = "Invalid credentials"
	MsgUsernameNotFound    = "Username tidak ditemukan"
	MsgPasswordWrong       = "Password salah"
	MsgTokenGenerationFail = "Gagal membuat token"
	MsgTokenInvalid        = "Token tidak valid"
	MsgAuthHeaderRequired  = "Authorization header wajib diisi"
)

// ========================
// USER MESSAGES
// ========================
const (
	MsgUserCreatedSuccess  = "User berhasil dibuat"
	MsgUserUpdateSuccess   = "User berhasil diperbarui"
	MsgUsernameExists      = "Username sudah digunakan"
	MsgUsernameRegistered  = "Username sudah terdaftar"
	MsgInvalidRole         = "Role tidak valid"
	MsgInvalidRoleValue    = "Role harus salah satu dari: admin, warehouse_admin, warehouse_manager"
	MsgUserNotFound        = "User tidak ditemukan"
	MsgPasswordEncryptFail = "Gagal mengenkripsi password"
	MsgUserSaveFail        = "Gagal menyimpan user"
	MsgUsersFetchSuccess   = "Daftar user berhasil didapatkan"
	MsgUserFetchSuccess    = "Data user berhasil didapatkan"
	MsgInvalidSession      = "Sesi tidak valid"
)

// ========================
// VALIDATION MESSAGES
// ========================
const (
	MsgValidationFailed   = "Validasi gagal"
	MsgValidationError    = "Validation error"
	MsgFieldRequired      = "Field ini wajib diisi"
	MsgUsernameMinLength  = "Username minimal 5 karakter"
	MsgPasswordMinLength  = "Password minimal 8 karakter"
	MsgInvalidFieldFormat = "Format tidak valid"
)

// ========================
// GENERAL ERROR MESSAGES
// ========================
const (
	MsgInternalServerError = "Terjadi kesalahan pada server"
	MsgInvalidRequest      = "Request payload tidak valid"
	MsgRenderFailed        = "Gagal memproses respon"
	MsgAdminOnlyAccess     = "Hanya admin yang dapat mengakses fitur ini"
	MsgUnauthorizedError   = "Unauthorized access"
	MsgForbiddenError      = "Forbidden access"
	MsgNotFoundError       = "Resource not found"
)

// ========================
// ITEM MESSAGES
// ========================
const (
	MsgItemCreatedSuccess   = "Barang berhasil dibuat"
	MsgItemCreatedFailed    = "Gagal menyimpan barang"
	MsgInvalidItemType      = "Jenis barang tidak valid"
	MsgInvalidUnit          = "Satuan barang tidak valid"
	MsgItemNotFound         = "Barang tidak ditemukan"
	MsgImageUploadFailed    = "Gagal mengupload gambar"
	MsgInvalidImageFormat   = "Format gambar tidak valid"
	MsgImageTooLarge        = "Ukuran gambar terlalu besar"
	MsgImageAllowedFormat   = "Hanya menerima format JPG, JPEG, atau PNG"
	MsgImageAllowedSizes    = "Maksimal ukuran gambar 5MB"
	MsgItemUpdatedSuccess   = "Item berhasil diperbarui"
	MsgItemUpdatedFailed    = "Gagal memperbarui item"
	MsgGetUpdatedItemFailed = "Gagal memuat data terupdate"
	MsgItemDeletedSuccess   = "Item berhasil dihapus"
	MsgItemDeleteFailed     = "Gagal menghapus item"
	MsgItemsFetchSuccess    = "Daftar item berhasil didapatkan"
	MsgItemFetchSuccess     = "Detail item berhasil didapatkan"
	MsgItemInUse            = "Item tidak dapat dihapus"
	MsgItemInUseDetail      = "Item sedang digunakan dalam %d transaksi"
)

// ========================
// ITEM TYPE
// ========================

const (
	MsgItemTypeCreatedSuccess = "Jenis barang berhasil dibuat"
	MsgItemTypeExists         = "Jenis barang sudah ada"
	MsgItemTypeExistsDetail   = "Nama jenis barang sudah terdaftar"
	MsgItemTypeCreateFailed   = "Gagal membuat jenis barang"
	MsgItemTypeNotFound       = "Jenis barang tidak ditemukan"
	MsgItemTypeUpdatedSuccess = "Jenis barang berhasil diperbarui"
	MsgItemTypeUpdateFailed   = "Gagal memperbarui jenis barang"
	MsgItemTypesFetchSuccess  = "Daftar jenis barang berhasil didapatkan"
	MsgItemTypeDeletedSuccess = "Jenis barang berhasil dihapus"
	MsgItemTypeDeleteFailed   = "Gagal menghapus jenis barang"
	MsgItemTypeInUse          = "Jenis barang tidak dapat dihapus"
	MsgItemTypeInUseDetail    = "Jenis barang sedang digunakan oleh  %d barang"
)

// ========================
// UNIT
// ========================

const (
	MsgUnitCreatedSuccess = "Satuan berhasil dibuat"
	MsgUnitUpdatedSuccess = "Satuan berhasil diperbarui"
	MsgUnitDeletedSuccess = "Satuan berhasil dihapus"
	MsgUnitExists         = "Satuan sudah ada"
	MsgUnitExistsDetail   = "Nama satuan sudah terdaftar"
	MsgUnitCreateFailed   = "Gagal membuat satuan"
	MsgUnitUpdateFailed   = "Gagal memperbarui satuan"
	MsgUnitDeleteFailed   = "Gagal menghapus satuan"
	MsgUnitNotFound       = "Satuan tidak ditemukan"
	MsgUnitsFetchSuccess  = "Daftar satuan berhasil didapatkan"
	MsgUnitFetchSuccess   = "Detail satuan berhasil didapatkan"
	MsgUnitInUseDetail    = "Satuan barang sedang digunakan oleh  %d barang"
)

// ========================
// TRANSACTION MESSAGES
// ========================

const (
	MsgTransactionCreatedSuccess = "Transaksi berhasil dibuat"
	MsgTransactionCreatedFailed  = "Gagal membuat transaksi"
	MsgTransactionFetchedSuccess = "Daftar transaksi berhasil didapatkan"
	MsgInsufficientStock         = "Stok tidak mencukupi"
	MsgTransactionDeletedSuccess = "Transaksi berhasil dihapus"
	MsgTransactionDeleteFailed   = "Gagal menghapus transaksi"
	MsgTransactionNotFound       = "Transaksi tidak ditemukan"
	MsgTransactionAdjustStock    = "Stok telah disesuaikan ke 0 karena penghapusan transaksi menyebabkan stok negatif"
)

// ========================
// ITEM REPORT MESSAGES
// ========================
const (
	MsgReportTitle              = "LAPORAN BARANG"
	MsgLowStockNote             = "*Menampilkan barang dengan stok di bawah minimum"
	MsgReportHeaderItemName     = "Nama Barang"
	MsgReportHeaderItemType     = "Jenis Barang"
	MsgReportHeaderUnit         = "Satuan"
	MsgReportHeaderCurrentStock = "Stok Saat Ini"
	MsgReportHeaderMinStock     = "Stok Minimum"
	MsgReportHeaderStatus       = "Status Stok"
	MsgStockStatusSafe          = "Aman"
	MsgStockStatusLow           = "Di Bawah Minimum"
	MsgFailedFetchItems         = "Gagal mengambil data barang"
	MsgFailedGenerateExcel      = "Gagal membuat file Excel"
	MsgReportFilenamePrefix     = "laporan_barang_"
	MsgReportTitleIn            = "LAPORAN BARANG MASUK"
	MsgReportTitleOut           = "LAPORAN BARANG KELUAR"
	MsgReportHeaderQuantity     = "Jumlah"
	MsgReportHeaderDate         = "Tanggal Transaksi"
	MsgReportHeaderDescription  = "Deskripsi"
	MsgReportFilenameTxPrefix   = "laporan_transaksi_"
	MsgInvalidDateRange         = "Range tanggal tidak valid"
	MsgFailedFetchTransactions  = "Gagal mengambil data transaksi"
	MsgInvalidTransactionType   = "Tipe transaksi tidak valid"
)

func GetReportTitleByType(txType string) string {
	switch txType {
	case TransactionTypeIn:
		return MsgReportTitleIn
	case TransactionTypeOut:
		return MsgReportTitleOut
	default:
		return ""
	}
}
