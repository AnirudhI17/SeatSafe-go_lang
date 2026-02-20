package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"ticketing/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.Database.URL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Migrations directory (run from backend/ directory)
	migrationsDir := "migrations"
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read migrations directory: %v\n", err)
		os.Exit(1)
	}

	// Collect all .up.sql files and sort them
	var migrationFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			migrationFiles = append(migrationFiles, entry.Name())
		}
	}
	sort.Strings(migrationFiles)

	fmt.Printf("Found %d migration files\n", len(migrationFiles))

	for _, filename := range migrationFiles {
		path := filepath.Join(migrationsDir, filename)
		sql, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Running migration: %s\n", filename)
		_, err = pool.Exec(ctx, string(sql))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Migration %s failed: %v\n", filename, err)
			os.Exit(1)
		}
		fmt.Printf("✓ %s completed\n", filename)
	}

	fmt.Println("\nAll migrations completed successfully!")
}
