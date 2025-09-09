package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pets_rest/internal/config"
	"pets_rest/internal/migrate"
)

func main() {
	var (
		up      = flag.Bool("up", false, "Run migrations")
		down    = flag.Bool("down", false, "Rollback migrations")
		version = flag.Bool("version", false, "Show current migration version")
	)
	flag.Parse()

	// Load configuration
	cfg := config.Load()

	// Get migrations path
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		log.Fatal("Failed to get migrations path:", err)
	}

	switch {
	case *up:
		if err := migrate.Run(cfg.DatabaseURL, migrationsPath); err != nil {
			log.Fatal("Failed to run migrations:", err)
		}
		fmt.Println("Migrations completed successfully")

	case *down:
		if err := migrate.Down(cfg.DatabaseURL, migrationsPath); err != nil {
			log.Fatal("Failed to rollback migrations:", err)
		}
		fmt.Println("Migrations rolled back successfully")

	case *version:
		v, dirty, err := migrate.Version(cfg.DatabaseURL, migrationsPath)
		if err != nil {
			log.Fatal("Failed to get migration version:", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %v)\n", v, dirty)

	default:
		fmt.Println("Usage: migrate -up|-down|-version")
		fmt.Println("  -up      Run migrations")
		fmt.Println("  -down    Rollback migrations")
		fmt.Println("  -version Show current migration version")
		os.Exit(1)
	}
}
