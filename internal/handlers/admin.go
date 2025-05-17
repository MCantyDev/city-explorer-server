package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/models"
	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

// All Admin-based Handlers
// Define what an Admin can do (Can Edit (Users table))
// Can Refresh column (call external API before Expiry date)
// See all data in the Server in the Frontend Admin Dashboard

func CheckAdminStatus(c *gin.Context) {
	isAdmin, exists := c.Get("isAdmin")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "isAdmin not found in context",
			"isAdmin": isAdmin,
		})
		return
	}
	if !isAdmin.(bool) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Admin status required",
			"isAdmin": isAdmin,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"error":   nil,
		"isAdmin": isAdmin,
	})
}

func GetUsers(c *gin.Context) {
	var users []models.User

	query := database.NewQueryBuilder("SELECT").Table("users").Build()
	_, err := database.Execute(&users, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured querying database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": users,
	})
}

func AddUser(c *gin.Context) {
	var req models.AddUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not hash password",
		})
		return
	}

	query := database.NewQueryBuilder("INSERT").Table("users").Columns("first_name", "last_name", "username", "email", "password", "is_admin").Values(6).Build()
	_, err = database.Execute(nil, query, req.FirstName, req.LastName, req.Username, req.Email, hashedPassword, req.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func EditUser(c *gin.Context) {
	var req models.EditUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Password != "" {
		// Hash Password
		hashedPassword, err := services.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not hash password",
			})
			return
		}

		query := database.NewQueryBuilder("UPDATE").Table("users").Columns("first_name", "last_name", "username", "email", "password", "is_admin").Where("id = ?").Build()
		_, err = database.Execute(nil, query, req.FirstName, req.LastName, req.Username, req.Email, hashedPassword, req.IsAdmin, req.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	} else {
		query := database.NewQueryBuilder("UPDATE").Table("users").Columns("first_name", "last_name", "username", "email", "is_admin").Where("id = ?").Build()
		_, err := database.Execute(nil, query, req.FirstName, req.LastName, req.Username, req.Email, req.IsAdmin, req.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func DeleteUser(c *gin.Context) {
	var req models.Delete

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not bind data to model in server",
		})
		return
	}

	query := database.NewQueryBuilder("DELETE").Table("users").Where("id = ?").Build()
	_, err := database.Execute(nil, query, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Removed ID - %d - from Users", req.Id),
	})
}

func GetCountries(c *gin.Context) {
	var countries []models.Country

	query := database.NewQueryBuilder("SELECT").Table("countries").Build()
	_, err := database.Execute(&countries, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured querying database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": countries,
	})
}

func GetCityWeatherTable(c *gin.Context) {
	var weatherReports []models.CityWeather

	query := database.NewQueryBuilder("SELECT").Table("city_weather").
		Columns("city_weather.id", "lat", "lon", "cities.name AS city_name", "countries.name AS country_name", "city_weather.data", "city_weather.created_at", "city_weather.updated_at", "city_weather.expiry_date").
		Join("JOIN cities ON cities.id=city_weather.city_id").Join("JOIN countries ON countries.id=city_weather.country_id").Build()
	_, err := database.Execute(&weatherReports, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured querying database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": weatherReports,
	})
}

func GetCitySightsTable(c *gin.Context) {
	var sights []models.CitySights

	query := database.NewQueryBuilder("SELECT").Table("city_sights").
		Columns("city_sights.id", "lat", "lon", "cities.name AS city_name", "countries.name AS country_name", "city_sights.data", "city_sights.created_at", "city_sights.updated_at", "city_sights.expiry_date").
		Join("JOIN cities ON cities.id=city_sights.city_id").Join("JOIN countries ON countries.id=city_sights.country_id").Build()
	_, err := database.Execute(&sights, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured querying database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": sights,
	})
}

func GetCityPoisTable(c *gin.Context) {
	var pois []models.CityPoi

	query := database.NewQueryBuilder("SELECT").Table("city_pois").
		Columns("city_pois.id", "cities.name AS city_name", "countries.name AS country_name", "xid", "city_pois.data", "city_pois.created_at", "city_pois.updated_at", "city_pois.expiry_date").
		Join("JOIN cities ON cities.id=city_pois.city_id").Join("JOIN countries ON countries.id=city_pois.country_id").Build()
	_, err := database.Execute(&pois, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error occured querying database",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": pois,
	})
}

func RefreshCountry(c *gin.Context) {
	countryCode := c.Query("country-code")
	if countryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'country-code' query parameter",
		})
		return
	}

	// Retreive Refreshed Data
	url := fmt.Sprintf(config.Cfg.RestCountriesAPI.URL, countryCode)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.RestCountriesAPI.Name + "' URL is not properly setup",
		})
		return
	}
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var validate models.RestCountriesRequest
	if err := json.Unmarshal(data, &validate); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid JSON from RestCountries: " + err.Error()})
		return
	}

	// UPDATE
	expiry := time.Now().AddDate(1, 0, 0)

	// This was strange...basically Rest Countries returns a json array [{}, {}...]
	// but im only being returned a single result (as im using Country codes)
	// So I am now just taking the FIRST element of each result and using that for Data...so the application doesnt get stuck on 'GetCountries'
	query := database.NewQueryBuilder("UPDATE").Table("countries").Columns("data", "expiry_date").Where("iso_code = ?").Build()
	_, err = database.Execute(nil, query, validate[0], expiry, countryCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func RefreshCityWeather(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' or 'lon' query parameter",
		})
		return
	}

	// Retreive Refreshed Data
	url := fmt.Sprintf(config.Cfg.OpenWeatherAPI.URL, lat, lon, config.Cfg.OpenWeatherAPI.Key)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenWeatherAPI.Name + "' URL is not set up properly",
		})
		return
	}
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// UPDATE
	expiry := time.Now().AddDate(0, 0, 1)

	query := database.NewQueryBuilder("UPDATE").Table("city_weather").Columns("data", "expiry_date").Where("ABS(lat - ?) < 0.0001 AND ABS(lon - ?) < 0.0001").Build()
	_, err = database.Execute(nil, query, data, expiry, lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func RefreshCitySights(c *gin.Context) {
	lat := c.Query("lat")
	lon := c.Query("lon")
	if lat == "" || lon == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' or 'lon' query parameter",
		})
		return
	}

	// Retreive Refreshed Data
	url := fmt.Sprintf(config.Cfg.OpenTripAPI.URL, lat, lon, config.Cfg.OpenTripAPI.Key)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenTripAPI.Name + "' URL is not setup properly",
		})
		return
	}
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// UPDATE
	expiry := time.Now().AddDate(0, 0, 1)

	query := database.NewQueryBuilder("UPDATE").Table("city_sights").Columns("data", "expiry_date").Where("ABS(lat - ?) < 0.0001 AND ABS(lon - ?) < 0.0001").Build()
	_, err = database.Execute(nil, query, data, expiry, lat, lon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func RefreshCityPoi(c *gin.Context) {
	xid := c.Query("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'xid' query parameter",
		})
		return
	}

	// Retreive Refreshed Data
	url := fmt.Sprintf(config.Cfg.OpenTripXIDAPI.URL, xid, config.Cfg.OpenTripAPI.Key)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenTripXIDAPI.Name + "' URL is not setup properly",
		})
		return
	}
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// UPDATE
	expiry := time.Now().AddDate(0, 0, 1)

	query := database.NewQueryBuilder("UPDATE").Table("city_pois").Columns("data", "expiry_date").Where("xid = ?").Build()
	_, err = database.Execute(nil, query, data, expiry, xid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func DeleteCountry(c *gin.Context) {
	var req models.Delete

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := database.NewQueryBuilder("DELETE").Table("countries").Where("id = ?").Build()
	_, err := database.Execute(nil, query, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func DeleteCityWeather(c *gin.Context) {
	var req models.Delete

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := database.NewQueryBuilder("DELETE").Table("city_weather").Where("id = ?").Build()
	_, err := database.Execute(nil, query, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func DeleteCitySights(c *gin.Context) {
	var req models.Delete

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := database.NewQueryBuilder("DELETE").Table("city_sights").Where("id = ?").Build()
	_, err := database.Execute(nil, query, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}

func DeleteCityPoi(c *gin.Context) {
	var req models.Delete

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	query := database.NewQueryBuilder("DELETE").Table("city_pois").Where("id = ?").Build()
	_, err := database.Execute(nil, query, req.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
