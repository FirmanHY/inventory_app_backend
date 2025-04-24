package main

import (
	"context"

	"inventory_app_backend/internal/config"
	"inventory_app_backend/internal/handlers"
	"inventory_app_backend/internal/routes"
	"inventory_app_backend/pkg/database"
	"inventory_app_backend/pkg/firebase"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load konfigurasi
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Inisialisasi database
	db, err := database.NewMySQLDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// inisialisasi firebase storage
	if err := firebase.InitializeStorage(); err != nil {
		log.Fatalf("Gagal inisialisasi storage: %v", err)
	}
	log.Println("Storage berhasil diinisialisasi")

	authHandler := &handlers.AuthHandler{DB: db}
	itemHandler := &handlers.ItemHandler{DB: db}
	itemTypeHandler := &handlers.ItemTypeHandler{DB: db}
	unitHandler := &handlers.UnitHandler{DB: db}
	transactionHandler := &handlers.TransactionHandler{DB: db}
	reportHandler := &handlers.ReportHandler{DB: db}
	summaryHandler := &handlers.SummaryHandler{DB: db}

	// Setup router
	router := routes.SetupRouter(
		authHandler,
		itemHandler,
		itemTypeHandler,
		unitHandler,
		transactionHandler,
		reportHandler,
		summaryHandler,
	)
	// Setup server
	port := config.Get("APP_PORT")
	if port == "" {
		port = "8080" // Default port
	}
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Jalankan server di goroutine
	go func() {
		log.Printf("Starting server on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Beri waktu 5 detik untuk selesaikan request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// Tutup koneksi database
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("Failed to get sql.DB: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Failed to close database: %v", err)
		}
	}

	log.Println("Server exited gracefully")
}
