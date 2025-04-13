package models

import (
	"time"
)

type User struct {
	UserID    string `gorm:"primaryKey;type:char(36)"`
	Username  string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	FullName  string
	Role      string `gorm:"type:ENUM('admin', 'warehouse_admin', 'warehouse_manager')"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
