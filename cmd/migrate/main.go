package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	var (
		direction = flag.String("direction", "up", "Migration direction: up, down, or force")
		steps     = flag.Int("steps", 0, "Number of migration steps (0 for all)")
		path      = flag.String("path", "", "Path to migration files")
		database  = flag.String("database", "", "Database URL")
		version   = flag.Int("version", 0, "Version for force command or targeted migration")
	)
	flag.Parse()

	// ---- Database URL ----
	databaseURL := *database
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		databaseURL = buildDatabaseURLFromEnv()
	}
	fmt.Println("Using database:", maskPassword(databaseURL))

	// ---- Migrations Path ----
	migrationsPath := *path
	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}
	if migrationsPath == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		if _, err := os.Stat("/app/migrations"); err == nil {
			migrationsPath = "/app/migrations"
		} else if _, err := os.Stat("/migrations"); err == nil {
			migrationsPath = "/migrations"
		} else {
			if filepath.Base(wd) == "migrate" {
				wd = filepath.Dir(filepath.Dir(wd))
			}
			migrationsPath = filepath.Join(wd, "internal", "migrations")
		}
	}
	migrationsURL := "file://" + migrationsPath
	fmt.Printf("Using migrations path: %s\n", migrationsURL)

	// ---- Create Migrate Instance ----
	m, err := migrate.New(migrationsURL, databaseURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer func() {
		srcErr, dbErr := m.Close()
		if srcErr != nil {
			log.Printf("Warning: source close error: %v", srcErr)
		}
		if dbErr != nil {
			log.Printf("Warning: database close error: %v", dbErr)
		}
	}()

	// ---- Handle Commands ----
	switch *direction {
	case "up":
		runUp(m, migrationsPath, *steps, *version)
	case "down":
		runDown(m, *steps)
	case "force":
		if *version == 0 {
			log.Fatalf("Version is required for force command")
		}
		if err := m.Force(*version); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		fmt.Printf("Forced migration version to %d successfully\n", *version)
	default:
		log.Fatalf("Invalid direction: %s. Use 'up', 'down', or 'force'", *direction)
	}
}

func buildDatabaseURLFromEnv() string {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if sslmode == "" {
		sslmode = "disable"
	}
	if host == "" || user == "" || pass == "" || name == "" {
		log.Fatalf("Missing required DB environment variables. Set DB_HOST, DB_USER, DB_PASSWORD, DB_NAME, DB_PORT, DB_SSLMODE or DATABASE_URL.")
	}

	url := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, name, sslmode,
	)
	return url
}

func maskPassword(url string) string {
	re := regexp.MustCompile(`(postgres:\/\/[^:]+:)[^@]+(@.*)`)
	return re.ReplaceAllString(url, "${1}****${2}")
}

func runUp(m *migrate.Migrate, migrationsPath string, steps int, targetVersion int) {
	current, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get current version: %v", err)
	}
	if dirty {
		log.Fatalf("Database is in a dirty state. Please fix manually or use --direction force.")
	}

	latest, err := getLatestMigrationVersion(migrationsPath)
	if err != nil {
		log.Fatalf("Failed to get latest migration version: %v", err)
	}

	fmt.Printf("Current DB version: %d | Latest available: %d\n", current, latest)

	if targetVersion > 0 {
		if targetVersion == int(current) {
			fmt.Printf("Database already at version %d\n", current)
			return
		}
		fmt.Printf("Migrating from version %d → %d\n", current, targetVersion)
		if err := m.Migrate(uint(targetVersion)); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to migrate to version %d: %v", targetVersion, err)
		}
		fmt.Println("Migration completed successfully.")
		return
	}

	if current == latest {
		fmt.Println("Database already up-to-date.")
		return
	}

	fmt.Printf("Migrating from version %d → %d\n", current, latest)
	if steps > 0 {
		err = m.Steps(steps)
	} else {
		err = m.Up()
	}
	if err != nil && err != migrate.ErrNoChange {
		if strings.Contains(err.Error(), "already exists") {
			fmt.Printf("Warning: Some objects already exist, continuing... %v\n", err)
		} else {
			log.Fatalf("Failed to run migration up: %v", err)
		}
	}
	fmt.Println("Migration up completed successfully.")
}

func runDown(m *migrate.Migrate, steps int) {
	current, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get current version: %v", err)
	}
	if dirty {
		log.Fatalf("Database is in a dirty state. Please fix manually or use --direction force.")
	}

	fmt.Printf("Current DB version: %d | Migrating down...\n", current)
	if steps > 0 {
		err = m.Steps(-steps)
	} else {
		err = m.Down()
	}
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migration down: %v", err)
	}
	fmt.Println("Migration down completed successfully.")
}

func getLatestMigrationVersion(path string) (uint, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return 0, err
	}
	re := regexp.MustCompile(`^(\d+)_`)
	var max uint
	for _, f := range files {
		if matches := re.FindStringSubmatch(f.Name()); len(matches) > 1 {
			v, _ := strconv.ParseUint(matches[1], 10, 64)
			if uint(v) > max {
				max = uint(v)
			}
		}
	}
	return max, nil
}
