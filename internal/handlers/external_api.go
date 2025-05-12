package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/MCantyDev/city-explorer-server/internal/models"
	"github.com/MCantyDev/city-explorer-server/internal/services"
	"github.com/gin-gonic/gin"
)

func GetCities(c *gin.Context) {
	city := c.Query("city")

	// Ensure City was set in the Get Request
	if city == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'city' query parameter",
		})
		return
	}

	// Encoding the Input for spaces and such
	encodedCity := url.QueryEscape(city)

	// Get the URL Template for Photon API
	url := fmt.Sprintf(config.Cfg.PhotonAPI.URL, encodedCity)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.PhotonAPI.Name + "' URL is not properly setup",
		})
		return
	}

	// Use City to Call External API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Map City Data
	var externalData models.PhotonRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

func GetCountry(c *gin.Context) {
	countryCode := c.Query("country-code")

	// Ensure Country Code was set in the Get Request
	if countryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'country-code' query parameter",
		})
		return
	}

	// Get the URL Template for Photon API
	url := fmt.Sprintf(config.Cfg.RestCountriesAPI.URL, countryCode)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.RestCountriesAPI.Name + "' URL is not properly setup",
		})
		return
	}

	// Use Country Code to Call External API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Map Country Data
	var externalData models.RestCountriesRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

// Weather and Such
func GetWeather(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")

	// Ensure Lat and Long are set in the Get Request
	if lat == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' query parameter",
		})
		return
	}
	if long == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'long' query parameter",
		})
		return
	}

	// Build URL String for OpenWeatherAPI
	url := fmt.Sprintf(config.Cfg.OpenWeatherAPI.URL, lat, long, config.Cfg.OpenWeatherAPI.Key)
	fmt.Println(url)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenWeatherAPI.Name + "' URL is not setup properly",
		})
		return
	}

	// Fetch Data from API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Unmarshal the Data into a usable form
	var externalData models.OpenWeatherRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

func GetTravelDestinations(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")

	// Ensure Lat and Long are set in the Get Request
	if lat == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' query parameter",
		})
		return
	}
	if long == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'long' query parameter",
		})
		return
	}

	// Build URL String for OpenWeatherAPI
	url := fmt.Sprintf(config.Cfg.OpenTripAPI.URL, lat, long, config.Cfg.OpenTripAPI.Key)
	fmt.Println(url)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenTripAPI.Name + "' URL is not setup properly",
		})
		return
	}

	// Fetch Data from API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Unmarshal the Data into a usable form
	var externalData models.OpenTripRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}
