package routes

import (
	"github.com/MCantyDev/city-explorer-server/internal/handlers"
	"github.com/MCantyDev/city-explorer-server/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// User Routes
	router.POST("/login", handlers.Login)
	router.POST("/sign-up", handlers.SignUp)

	// Grouping
	auth := router.Group("/auth")
	auth.Use(middleware.SessionAuthMiddleware())
	{
		auth.GET("/profile", handlers.GetProfile)
		auth.GET("/logout", handlers.Logout)
		auth.GET("/get-country", handlers.GetCountry)
		auth.GET("/get-cities", handlers.GetCities)
		auth.GET("/get-city-weather", handlers.GetWeather)
		auth.GET("/get-city-sights", handlers.GetTravelDestinations)
		auth.GET("/get-city-poi", handlers.GetTravelDestination)
		// auth.GET("/check-admin-status", handlers.CheckAdminStatus)
	}

	// Admin group
	admin := router.Group("/admin")
	admin.Use(middleware.SessionAuthMiddleware(), middleware.AdminMiddleware())
	{
		admin.GET("/get-users", handlers.GetUsers)
		admin.GET("/get-countries", handlers.GetCountries)
		admin.GET("/get-city-weather", handlers.GetCityWeatherTable)
		admin.GET("/get-city-sights", handlers.GetCitySightsTable)
		admin.GET("/get-city-pois", handlers.GetCityPoisTable)
		admin.POST("/add-user", handlers.AddUser)
		admin.PATCH("/edit-user", handlers.EditUser)
		admin.PATCH("/refresh-country", handlers.RefreshCountry)
		admin.PATCH("/refresh-city-weather", handlers.RefreshCityWeather)
		admin.PATCH("/refresh-city-sights", handlers.RefreshCitySights)
		admin.PATCH("/refresh-city-poi", handlers.RefreshCityPoi)
		admin.DELETE("/delete-user", handlers.DeleteUser)
		admin.DELETE("/delete-country", handlers.DeleteCountry)
		admin.DELETE("/delete-city-weather", handlers.DeleteCityWeather)
		admin.DELETE("/delete-city-sights", handlers.DeleteCitySights)
		admin.DELETE("/delete-city-poi", handlers.DeleteCityPoi)
	}
}
