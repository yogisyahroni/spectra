package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"spectra-backend/internal/config"
	"spectra-backend/internal/database"
	"spectra-backend/internal/routes"
)

func main() {
	log.Println("🚀 Starting SPECTRA Backend Server...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Failed to load configuration: %v", err)
	}

	log.Printf("📋 Environment: %s", cfg.Environment)
	log.Printf("🔗 Database: %s", maskDatabaseURL(cfg.GetDatabaseURL()))

	// Connect to database
	log.Println("🔌 Connecting to database...")
	db, err := database.New(cfg.GetDatabaseURL())
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Database connected successfully")

	// Run migrations
	log.Println("📦 Running database migrations...")
	if err := database.RunMigrations(db.Pool); err != nil {
		log.Fatalf("❌ Failed to run migrations: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

	// Setup routes
	handler := routes.SetupRoutes(db.Pool)

	// Create server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🌐 Server listening on http://%s", addr)
		log.Println("📍 API Endpoints:")
		log.Println("   GET  /api/health          - Health check")
		log.Println("   GET  /api/nodes           - List nodes")
		log.Println("   POST /api/nodes           - Create node")
		log.Println("   GET  /api/nodes/nearby    - Get nearby nodes")
		log.Println("   GET  /api/geojson/nodes   - Get nodes as GeoJSON")
		log.Println("   GET  /api/cables          - List cables")
		log.Println("   GET  /api/customers       - List customers")
		log.Println("   GET  /api/connections     - List connections")
		log.Println("")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server forced to shutdown: %v", err)
	}

	log.Println("👋 Server exited gracefully")
}

// maskDatabaseURL masks sensitive information in the database URL
func maskDatabaseURL(url string) string {
	// Simple masking - hide password
	// postgresql://user:password@host:port/db -> postgresql://user:****@host:port/db
	if len(url) > 20 {
		return url[:15] + "****" + url[len(url)-20:]
	}
	return "****"
}
