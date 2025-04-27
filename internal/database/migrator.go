package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

/*

Purpose of Migrator
-----------------

- Handles the Initialisation of Database (if DB does not exist)
- Handles Migrations

*/

func migrate(database *gorm.DB) error {
	// Path to the Migrations folder
	migrationsPath := "./internal/database/migrations"

	files, err := readMigrationDirectory(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fullSQL, err := readMigrationFile(file, migrationsPath)
		if err != nil {
			return err
		}

		splitSQL := splitMigrationQueries(fullSQL)
		err = executeMigrationQueries(splitSQL, database)
		if err != nil {
			return err
		}
	}
	return nil // Success!
}

// Utility Functions
func readMigrationDirectory(migrationsPath string) ([]os.DirEntry, error) {
	files, err := os.ReadDir(migrationsPath)

	if err != nil {
		return nil, fmt.Errorf("error occured while reading migration directory: %s", err)
	}
	return files, nil
}

func readMigrationFile(file os.DirEntry, migrationsPath string) (string, error) {
	if strings.HasSuffix(file.Name(), ".sql") {
		filePath := filepath.Join(migrationsPath, file.Name())
		fmt.Printf("Applying Migration: %s\n", file.Name())

		// Read the file
		migrationSQL, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(migrationSQL), nil // Success!
	}
	return "", fmt.Errorf("failed to read migration file '%s': 'Must end with .sql'", file.Name())
}

func splitMigrationQueries(queryString string) []string {
	queries := strings.Split(queryString, ";")

	for i := range queries {
		queries[i] = strings.TrimSpace(queries[i])
	}

	return queries
}

func executeMigrationQueries(queries []string, database *gorm.DB) error {
	for _, query := range queries {
		if query == "" {
			continue
		}

		if err := database.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to Execute Query: %s", err)
		}
	}
	return nil // Success!
}
