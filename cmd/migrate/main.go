package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		direction = flag.String("direction", "up", "Migration direction: up or down")
		steps     = flag.Int("steps", 0, "Number of migration steps (0 for all)")
	)
	flag.Parse()

	// Get database URL from environment or use default
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://booking_user:booking_pass@localhost:5432/booking_db?sslmode=disable"
	}

	// Get the absolute path to migrations directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	
	var migrationsPath string
	
	// Check if we're in a container environment
	if _, err := os.Stat("/migrations"); err == nil {
		// We're in a container with /migrations mounted
		migrationsPath = "/migrations"
	} else {
		// We're running locally
		if filepath.Base(wd) == "migrate" {
			wd = filepath.Dir(filepath.Dir(wd))
		}
		migrationsPath = filepath.Join(wd, "internal", "migrations")
	}
	
	migrationsURL := "file://" + migrationsPath
	fmt.Printf("Using migrations path: %s\n", migrationsURL)

	// Create migrate instance
	m, err := migrate.New(
		migrationsURL,
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	// Run migration based on direction
	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migration up: %v", err)
		}
		fmt.Println("Migration up completed successfully")

	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migration down: %v", err)
		}
		fmt.Println("Migration down completed successfully")

	default:
		log.Fatalf("Invalid direction: %s. Use 'up' or 'down'", *direction)
	}
}
