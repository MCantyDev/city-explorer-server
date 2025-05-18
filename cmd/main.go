package main

import (
	"log"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/routes"
	"github.com/gin-contrib/cors"
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

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(router)
	router.Run(":5050")

}
