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
	router.POST("/sign-up", handlers.SignUp)

	// Grouping
	auth := router.Group("/auth")
	auth.Use(middleware.RefreshSessionTokenMiddleware(), middleware.CookieMiddleware())
	{
		// Used to get User profile information, userful for User Page
		auth.GET("/profile", handlers.GetProfile)

		auth.GET("/get-country", handlers.GetCountry) // Rest-Countries API
		auth.GET("/get-cities", handlers.GetCities)   // PhotonAPI

		// City Specific Endpoints
		auth.GET("/get-city-weather", handlers.GetWeather)           // OpenWeatherMap API
		auth.GET("/get-city-sights", handlers.GetTravelDestinations) // OpenTripMap API
		auth.GET("/get-city-poi", handlers.GetTravelDestination)

		// Check Admin Status (Here as any user can technically check status)
		auth.GET("/check-admin-status", handlers.CheckAdminStatus)
	}

	// Admin group
	admin := router.Group("/admin")
	admin.Use(middleware.RefreshSessionTokenMiddleware(), middleware.CookieMiddleware(), middleware.AdminMiddleware())
	{
		// Get Table Data
		admin.GET("/get-users", handlers.GetUsers)
		admin.GET("/get-countries", handlers.GetCountries)
		admin.GET("/get-city-weather", handlers.GetCityWeatherTable)
		admin.GET("/get-city-sights", handlers.GetCitySightsTable)
		admin.GET("/get-city-pois", handlers.GetCityPoisTable)

		// Edit Users (Only Selectable Editable Data - The rest just Query the external APIs again)
		admin.POST("/add-user", handlers.AddUser)
		admin.POST("/edit-user", handlers.EditUser)
		admin.POST("/delete-user", handlers.DeleteUser)

		// Refresh (Using GET because...why not)
		admin.GET("/refresh-country", handlers.RefreshCountry)
		admin.GET("/refresh-city-weather", handlers.RefreshCityWeather)
		admin.GET("/refresh-city-sights", handlers.RefreshCitySights)
		admin.GET("/refresh-city-poi", handlers.RefreshCityPoi)

		// Delete
		admin.POST("/delete-country", handlers.DeleteCountry)
		admin.POST("/delete-city-weather", handlers.DeleteCityWeather)
		admin.POST("/delete-city-sights", handlers.DeleteCitySights)
		admin.POST("/delete-city-poi", handlers.DeleteCityPoi)
	}
}
