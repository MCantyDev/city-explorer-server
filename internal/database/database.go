package database

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gorm.io/gorm"
)

var DB *gorm.DB = nil

/*

Plan:
- Connect to the MYSQL Server (No Database)
- See if the "city_explorer_app" database exists
  - if exists -> Pass
  - if not exists -> run initialisation
- Reconnect WITH the "city_explorer_app" database
- Run migrations

*/

func Connect(dbName string) error {
	// Connect to MySQL Server
	serverConnection, err := connectToServer()
	if err != nil {
		return err
	}

	result, err := checkDatabaseExists(serverConnection, dbName)
	if err != nil {
		return err
	}

	var firstPass bool = false
	if !result {
		initialiseDatabase(serverConnection, dbName)
		firstPass = true

		err := close(serverConnection)
		if err != nil {
			return err
		}
	}

	// Connect to Database
	DB, err = connectToDatabase(dbName)
	if err != nil {
		return err
	}

	migrate(DB, firstPass)
	return nil // Success!
}

func Execute(model any, query string, args ...any) (any, error) {
	queryUpper := strings.ToUpper(strings.TrimSpace(query))

	switch {
	case strings.HasPrefix(queryUpper, "SELECT"):
		result := DB.Raw(query, args...).Scan(model)
		return model, result.Error

	case strings.HasPrefix(queryUpper, "INSERT") && model != nil:
		err := DB.Create(model).Error
		return model, err

	default:
		res := DB.Exec(query, args...)
		return res.RowsAffected, res.Error
	}
}

func initialiseDatabase(server *gorm.DB, dbName string) {
	queryBytes, err := os.ReadFile("./internal/database/migrations/initialisation/000_initialisation.sql")
	if err != nil {
		log.Fatal("Error Reading Migration Initialisation File: ", err)
	}

	query := strings.ReplaceAll(string(queryBytes), "{{SQL_DATABASE}}", dbName)
	if err := server.Exec(query).Error; err != nil {
		log.Fatal("Error Initialising Database: ", err)
	}
	fmt.Println("Successfully Initialised Database")
}
