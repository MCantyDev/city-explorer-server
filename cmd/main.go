package main

import (
	"log"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load environment variables
	config.Load()

	// Connect to the Database
	err := database.Connect(config.Cfg.Database.Name)
	if err != nil {
		log.Fatalf("Error Occured: %s", err)
	}

	// Setup Gin Server Router
	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.SetupRoutes(router)
	router.Run(":5050")
}
