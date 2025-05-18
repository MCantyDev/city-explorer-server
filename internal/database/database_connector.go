package database

import (
	"errors"
	"fmt"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

/*

Purpose of Database Connector
-----------------

- Handles connecting to Database (ensure Database exists in Server)
- Allow easy closing of Database (in case required)

*/

// connectToServer - establishes a connection to the MySQL server (Without selecting a DB) -> Allows a Database check to be ran to see if we need to initialise Database
func connectToServer() (*gorm.DB, error) {
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=true&loc=Local",
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
	)
	conn, err := gorm.Open(mysql.Open(dsnWithoutDB), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL Server: %s", err)
	}
	return conn, nil
}

// connectToDatabase - Establishes a connection to the MySQL Database

func connectToDatabase(dbName string) (*gorm.DB, error) {
	dsnWithDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		config.Cfg.Database.User,
		config.Cfg.Database.Password,
		config.Cfg.Database.Host,
		config.Cfg.Database.Port,
		dbName,
	)

	conn, err := gorm.Open(mysql.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL Database: %s", err)
	}
	return conn, nil
}

// checkDatabaseExist - Check if a Database exists within a MySQL Server
func checkDatabaseExists(server *gorm.DB, dbName string) (bool, error) {
	var result string
	query := fmt.Sprintf("SHOW DATABASES LIKE '%s';", dbName)

	dbResult := server.Raw(query).Scan(&result)
	if dbResult.Error != nil {
		if errors.Is(dbResult.Error, gorm.ErrRecordNotFound) {
			return false, dbResult.Error
		}
	}
	return result == dbName, nil
}

// close - Close a gorm.DB database using the underlying sqlDB struct
func close(db *gorm.DB) error {
	sqlDB, err := db.DB()

	if err != nil {
		return fmt.Errorf("error occured while closing connection to database: %s", err)
	}
	return sqlDB.Close()
}
