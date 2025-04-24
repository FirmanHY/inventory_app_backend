package database

import (
	"fmt"
	"inventory_app_backend/internal/config"

	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLDB() (*gorm.DB, error) {
	// Validasi environment variables
	requiredConfigs := []string{"DB_USER", "DB_HOST", "DB_PORT", "DB_NAME"}
	for _, cfg := range requiredConfigs {
		if config.Get(cfg) == "" {
			return nil, fmt.Errorf("missing required config: %s", cfg)
		}
	}

	// Format DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Get("DB_USER"),
		config.Get("DB_PASSWORD"), // <- ini yang kurang tadi
		config.Get("DB_HOST"),
		config.Get("DB_PORT"),
		config.Get("DB_NAME"),
	)

	// Retry mechanism
	var db *gorm.DB
	var err error
	maxRetries := 5
	retryDelay := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // Hanya untuk migrasi
		})
		if err == nil {
			break
		}
		log.Printf("Failed to connect to database (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Atur connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Konfigurasi connection pool
	sqlDB.SetMaxOpenConns(100)          // Maksimum koneksi terbuka
	sqlDB.SetMaxIdleConns(10)           // Maksimum koneksi idle
	sqlDB.SetConnMaxLifetime(time.Hour) // Maksimum umur koneksi

	// Test koneksi
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database")
	return db, nil
}
