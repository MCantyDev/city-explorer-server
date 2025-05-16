package database

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

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

// Rebuilt to be more robust and function with the newly added Triggers (Pain in the ass)
func splitMigrationQueries(queryString string) []string {
	var queries []string
	var currentQuery strings.Builder
	var prevCh rune

	inSingleQuote := false
	inDoubleQuote := false
	inLineComment := false
	inBlockComment := false
	beginEndDepth := 0

	var window []rune // Sliding window for detecting ' begin ' and ' end '

	for _, ch := range queryString {
		// Handle quotes
		if !inDoubleQuote && !inLineComment && !inBlockComment && ch == '\'' {
			inSingleQuote = !inSingleQuote
		} else if !inSingleQuote && !inLineComment && !inBlockComment && ch == '"' {
			inDoubleQuote = !inDoubleQuote
		}
		// Handle line comments (--) -> Cannot be Flipped like quotes as '' and "" are repeatable to end them...while ---- does not end the line comment
		if !inSingleQuote && !inDoubleQuote && !inBlockComment && prevCh == '-' && ch == '-' {
			inLineComment = true
		}
		if inLineComment && ch == '\n' {
			inLineComment = false
		}
		// Handle block comments (/* */) -> Similarly Block comments need a /* and */ doing /* /* or */ */ does not start and end comment
		if !inSingleQuote && !inDoubleQuote && !inLineComment && prevCh == '/' && ch == '*' {
			inBlockComment = true
		}
		if inBlockComment && prevCh == '*' && ch == '/' {
			inBlockComment = false
		}

		// Sliding Window
		window = append(window, unicode.ToLower(ch)) // Add rune to window
		if len(window) > 5 {                         // if window is 5 chars long
			window = window[1:] // remove the first index (by slicing from the second index)
		}
		windowStr := string(window) // Convert to string

		// Detect BEGIN / END depth changes
		if !inSingleQuote && !inDoubleQuote && !inLineComment && !inBlockComment {
			if windowStr == "begin" {
				beginEndDepth++
			} else if beginEndDepth > 0 && windowStr == "end" {
				beginEndDepth--
			}
		}

		// Add current character to builder
		currentQuery.WriteRune(ch)

		// Split Statement if ; is found outside all checks
		if !inSingleQuote && !inDoubleQuote && !inLineComment && !inBlockComment && beginEndDepth == 0 && ch == ';' {
			queries = append(queries, strings.TrimSpace(currentQuery.String()))
			currentQuery.Reset()
		}

		prevCh = ch
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
