package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

/*

Purpose of Migrator
-----------------

- Handles the Initialisation of Database (if DB does not exist)
- Handles Migrations

*/

// Migration with Check -> General Purpose, Ran at the Start of Connection no matter what
func migrate(database *gorm.DB, firstPass bool) error {
	// Path to the Migrations folder
	migrationsPath := "./internal/database/migrations"

	// Get Current Migration - Default to 0 in case of First Pass (True if Database being Initialised)
	var migrationInt int = 0
	if !firstPass {
		query := NewQueryBuilder("SELECT").Table("migration_tracker").Columns("migration_int").Build()
		_, err := Execute(&migrationInt, query)
		if err != nil {
			return err
		}
	}

	files, err := readMigrationDirectory(migrationsPath)
	if err != nil {
		return err
	}

	var highestMigration int = 0

	for _, file := range files {

		number, err := splitFileNames(file.Name())
		if err != nil {
			return err
		}

		if number > migrationInt {
			fullSQL, err := readMigrationFile(file, migrationsPath)
			if err != nil {
				return err
			}
			fmt.Printf("Applying Migration: %s\n", file.Name())

			splitSQL := splitMigrationQueries(fullSQL)
			err = executeMigrationQueries(splitSQL, database)
			if err != nil {
				return err
			}
			highestMigration = number
		}
	}

	if highestMigration > 0 {
		fmt.Println("Executing the update query")
		query := NewQueryBuilder("UPDATE").Table("migration_tracker").Columns("migration_int").Where("id = 1").Build()
		_, err = Execute(nil, query, highestMigration)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("No migration was applied - Database up to date")
	}

	return nil // Success!
}

// Utility Functions
func readMigrationDirectory(migrationsPath string) ([]os.DirEntry, error) {
	allFiles, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("error occured while reading migration directory: %s", err)
	}

	var validFiles []os.DirEntry
	for _, entry := range allFiles {
		if entry.IsDir() {
			continue
		}

		if !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}

		parts := strings.SplitN(entry.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}

		_, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}

		validFiles = append(validFiles, entry)
	}
	return validFiles, nil
}

func readMigrationFile(file os.DirEntry, migrationsPath string) (string, error) {
	if strings.HasSuffix(file.Name(), ".sql") {
		filePath := filepath.Join(migrationsPath, file.Name())

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

func splitFileNames(filename string) (int, error) {
	// Split the filename at the first "_" in 2 parts
	parts := strings.SplitN(filename, "_", 2)

	number, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, err
	}

	return number, nil
}
