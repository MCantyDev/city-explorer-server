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
		// Used to get User profile information, userful for User Page
		auth.GET("/profile", handlers.GetProfile)

		auth.GET("/get-country", handlers.GetCountry) // Rest-Countries API
		auth.GET("/get-cities", handlers.GetCities)   // PhotonAPI

		// City Specific Endpoints
		auth.GET("/get-city-weather", handlers.GetWeather)           // OpenWeatherMap API
		auth.GET("/get-city-sights", handlers.GetTravelDestinations) // OpenTripMap API
		auth.GET("/get-city-poi", handlers.GetTravelDestination)
	}
}
