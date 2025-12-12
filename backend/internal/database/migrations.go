package database

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

// RunMigrations executes all SQL migration files in order
func RunMigrations(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create migrations tracking table
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			filename VARCHAR(255) UNIQUE NOT NULL,
			executed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migration files
	entries, err := fs.ReadDir(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Sort by filename
	var filenames []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			filenames = append(filenames, entry.Name())
		}
	}
	sort.Strings(filenames)

	// Execute each migration
	for _, filename := range filenames {
		// Check if already executed
		var count int
		err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM schema_migrations WHERE filename = $1", filename).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status for %s: %w", filename, err)
		}

		if count > 0 {
			log.Printf("Migration %s already executed, skipping...", filename)
			continue
		}

		// Read and execute migration
		content, err := fs.ReadFile(migrationsFS, "migrations/"+filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", filename, err)
		}

		log.Printf("Executing migration: %s", filename)

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", filename, err)
		}

		_, err = tx.Exec(ctx, string(content))
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to execute migration %s: %w", filename, err)
		}

		// Record migration
		_, err = tx.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", filename)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("failed to record migration %s: %w", filename, err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", filename, err)
		}

		log.Printf("Migration %s completed successfully", filename)
	}

	return nil
}
