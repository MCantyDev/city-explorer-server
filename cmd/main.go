package main

import (
	"log"
	"os"

	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/routes"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	// Read the ENV file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}

	database.Connect(os.Getenv("DB_NAME"))

	// Test Cases for the Database Querying
	// query := database.NewQueryBuilder("INSERT").Table("users").Columns("username", "password").Values(2)
	// database.Execute(query.Build(), "MarkC", "saddadasdad")

	// query := database.NewQueryBuilder("UPDATE").Table("users").Columns("password").Where("id = 1")
	// database.Execute(query.Build(), "scoobydoo")

	// query := database.NewQueryBuilder("DELETE").Table("users").Where("id = 1")
	// database.Execute(query.Build())

	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":5050")
}
