package models

import (
	"time"
)

type Unit struct {
	UnitID    string `gorm:"primaryKey;type:char(36)"`
	UnitName  string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
