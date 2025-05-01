package main

import (
	"log"
	"os"

	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/routes"
	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
)

func main() {
	// Read the ENV file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Occured: could not read .env file (Ensure it is setup in root directory)")
	}

	// Connect to the Database
	err = database.Connect(os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("Error Occured: %s", err)
	}

	// Setup JWT Token system
	services.InitJWT()
	if err != nil {
		log.Fatalf("Error Occured: %s", err)
	}

	// Setup Gin Server Router
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":5050")
}
