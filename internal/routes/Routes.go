package routes

import (
	"github.com/MCantyDev/city-explorer-server/internal/handlers"
	"github.com/MCantyDev/city-explorer-server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// Setup the Router Routes

	// User Routes
	router.POST("/login", handlers.Login)
	router.POST("/signup", handlers.SignUp)

	// Grouping
	auth := router.Group("/auth")
	auth.Use(middleware.JWTMiddleware())
	{
		auth.GET("/profile", handlers.GetProfile)
	}
}
