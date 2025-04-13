package dto

import "time"

type ItemReportDTO struct {
	ItemName     string
	TypeName     string
	UnitName     string
	Stock        int
	MinimumStock int
	Status       string
}

type TransactionReportDTO struct {
	ItemName    string
	TypeName    string
	Quantity    int
	Date        time.Time
	Description string
	Type        string // in/out
}
