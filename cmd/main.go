package main

import (
	"github.com/MCantyDev/city-explorer-server/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	routes.SetupRoutes(router)

	router.Run(":5050")
}
