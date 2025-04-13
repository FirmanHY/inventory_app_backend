package models

import (
	"time"
)

type ItemType struct {
	TypeID    string `gorm:"primaryKey;type:char(36)"`
	TypeName  string `gorm:"unique;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
