package models

import (
	"time"
)

type Transaction struct {
	TransactionID   string    `gorm:"primaryKey;type:char(36)"`
	ItemID          string    `gorm:"type:char(36);not null"`
	Date            time.Time `gorm:"type:date;not null"`
	Quantity        int       `gorm:"not null"`
	TransactionType string    `gorm:"type:ENUM('in', 'out');not null"`
	Description     string
	UserID          string `gorm:"type:char(36);not null"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relations
	Item Item `gorm:"foreignKey:ItemID;references:ItemID"`
	User User `gorm:"foreignKey:UserID;references:UserID"`
}
