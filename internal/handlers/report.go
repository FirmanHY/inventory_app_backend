// internal/handlers/report.go
package handlers

import (
	constants "inventory_app_backend/internal/constant"
	"inventory_app_backend/internal/dto"
	"inventory_app_backend/internal/models"
	"inventory_app_backend/internal/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ReportHandler struct {
	DB *gorm.DB
}

func (h *ReportHandler) GenerateItemReport(c *gin.Context) {
	lowStockOnly, _ := strconv.ParseBool(c.DefaultQuery("low_stock_only", "false"))

	query := h.DB.Preload("Type").Preload("Unit")
	if lowStockOnly {
		query = query.Where("stock < minimum_stock")
	}

	var items []models.Item
	if err := query.Find(&items).Error; err != nil {
		utils.ServerError(c, constants.MsgFailedFetchItems, err)
		return
	}

	var reportData []dto.ItemReportDTO
	for _, item := range items {
		status := constants.MsgStockStatusSafe
		if item.Stock < item.MinimumStock {
			status = constants.MsgStockStatusLow
		}

		reportData = append(reportData, dto.ItemReportDTO{
			ItemName:     item.ItemName,
			TypeName:     item.Type.TypeName,
			UnitName:     item.Unit.UnitName,
			Stock:        item.Stock,
			MinimumStock: item.MinimumStock,
			Status:       status,
		})
	}

	f := excelize.NewFile()
	defer f.Close()

	sheetName := constants.MsgReportTitle
	if lowStockOnly {
		sheetName = constants.MsgLowStockNote
	}

	f.SetCellValue(sheetName, "A1", constants.MsgReportTitle)
	if lowStockOnly {
		f.SetCellValue(sheetName, "A2", constants.MsgLowStockNote)
	}
	f.NewSheet(sheetName)
	f.DeleteSheet("Sheet1")

	headers := []string{
		constants.MsgReportHeaderItemName,
		constants.MsgReportHeaderItemType,
		constants.MsgReportHeaderUnit,
		constants.MsgReportHeaderCurrentStock,
		constants.MsgReportHeaderMinStock,
		constants.MsgReportHeaderStatus,
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
		f.SetCellStyle(sheetName, cell, cell, h.headerStyle(f))
	}

	for rowIdx, data := range reportData {
		row := rowIdx + 2
		values := []interface{}{
			data.ItemName,
			data.TypeName,
			data.UnitName,
			data.Stock,
			data.MinimumStock,
			data.Status,
		}

		for colIdx, value := range values {
			cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
			f.SetCellValue(sheetName, cell, value)
		}
	}

	for i := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(sheetName, col, col, 20)
	}

	filename := constants.MsgReportFilenamePrefix + time.Now().Format("20060102_150405") + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		utils.ServerError(c, constants.MsgFailedGenerateExcel, err)
		return
	}

	c.Abort()
}

func (h *ReportHandler) headerStyle(f *excelize.File) int {
	style, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#DFF0D8"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	return style
}

func (h *ReportHandler) GenerateTransactionReport(c *gin.Context) {
	// Parse query params
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))
	txType := c.Query("type")

	// Validasi transaction type
	if txType != "" && txType != constants.TransactionTypeIn && txType != constants.TransactionTypeOut {
		utils.BadRequest(c, constants.MsgInvalidTransactionType, nil)
		return
	}

	// Validasi date range
	if !startDate.IsZero() && !endDate.IsZero() && startDate.After(endDate) {
		utils.BadRequest(c, constants.MsgInvalidDateRange, nil)
		return
	}

	// Build query
	query := h.DB.Preload("Item.Type").Preload("Item").Preload("User").
		Model(&models.Transaction{}).
		Order("date ASC")

	// Apply filters
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("date BETWEEN ? AND ?", startDate.Format("2006-01-02"),
			endDate.Format("2006-01-02"))
	}
	if txType != "" {
		query = query.Where("transaction_type = ?", txType)
	}

	var transactions []models.Transaction
	if err := query.Find(&transactions).Error; err != nil {
		utils.ServerError(c, constants.MsgFailedFetchTransactions, err)
		return
	}

	// Convert to DTO
	var reportData []dto.TransactionReportDTO
	for _, tx := range transactions {
		reportData = append(reportData, dto.TransactionReportDTO{
			ItemName:    tx.Item.ItemName,
			TypeName:    tx.Item.Type.TypeName,
			Quantity:    tx.Quantity,
			Date:        tx.Date,
			Description: tx.Description,
			Type:        tx.TransactionType,
		})
	}

	// Buat file Excel
	f := excelize.NewFile()
	defer f.Close()

	// Tentukan sheet yang akan dibuat
	sheets := make(map[string]string)
	if txType == "" {
		// Jika tidak ada filter type, buat kedua sheet
		sheets[constants.TransactionTypeIn] = constants.MsgReportTitleIn
		sheets[constants.TransactionTypeOut] = constants.MsgReportTitleOut
	} else {
		// Jika ada filter type, buat sheet sesuai type
		sheets[txType] = constants.GetReportTitleByType(txType)
	}

	headers := []string{
		constants.MsgReportHeaderItemName,
		constants.MsgReportHeaderItemType,
		constants.MsgReportHeaderQuantity,
		constants.MsgReportHeaderDate,
		constants.MsgReportHeaderDescription,
	}

	// Buat sheet dan isi data
	for sheetType, sheetName := range sheets {
		index, _ := f.NewSheet(sheetName)
		f.SetActiveSheet(index)

		// Set header
		for i, header := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheetName, cell, header)
			f.SetCellStyle(sheetName, cell, cell, h.headerStyle(f))
		}

		// Filter data per type
		var filteredData []dto.TransactionReportDTO
		for _, data := range reportData {
			if data.Type == sheetType {
				filteredData = append(filteredData, data)
			}
		}

		// Isi data
		for rowIdx, data := range filteredData {
			row := rowIdx + 2
			values := []interface{}{
				data.ItemName,
				data.TypeName,
				data.Quantity,
				data.Date.Format("2006-01-02"),
				data.Description,
			}

			for colIdx, value := range values {
				cell, _ := excelize.CoordinatesToCellName(colIdx+1, row)
				f.SetCellValue(sheetName, cell, value)
			}
		}

		// Set column width
		for i := range headers {
			col, _ := excelize.ColumnNumberToName(i + 1)
			f.SetColWidth(sheetName, col, col, 20)
		}
	}

	// Hapus sheet default jika ada sheet lain yang dibuat
	if len(sheets) > 0 {
		f.DeleteSheet("Sheet1")
	}

	// Set headers untuk download
	filename := constants.MsgReportFilenameTxPrefix + time.Now().Format("20060102_150405") + ".xlsx"
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	if err := f.Write(c.Writer); err != nil {
		utils.ServerError(c, constants.MsgFailedGenerateExcel, err)
		return
	}

	c.Abort()
}
