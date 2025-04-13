package models

import (
	"time"
)

type Item struct {
	ItemID       string `gorm:"primaryKey;type:char(36)"`
	TypeID       string `gorm:"type:char(36);not null"`
	UnitID       string `gorm:"type:char(36);not null"`
	ItemName     string `gorm:"not null"`
	Stock        int    `gorm:"not null;default:0"`
	MinimumStock int    `gorm:"not null;default:0"`
	Image        string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	// Relations
	Type ItemType `gorm:"foreignKey:TypeID;references:TypeID"`
	Unit Unit     `gorm:"foreignKey:UnitID;references:UnitID"`
}
