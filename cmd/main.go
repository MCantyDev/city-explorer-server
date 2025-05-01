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

	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":5050")
}
