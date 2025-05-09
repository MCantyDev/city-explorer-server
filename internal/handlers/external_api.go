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

func GetCity(c *gin.Context) {
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
	url := fmt.Sprintf(config.Cfg.PhotonAPI.URL+"?q=%s", encodedCity)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'PHOTON_URL' is not properly setup",
		})
		return
	}

	// Use City to Call External API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Map City Data
	var externalData models.PhotonRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
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
	url := fmt.Sprintf(config.Cfg.RestCountriesAPI.URL+"%s", countryCode)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'REST_COUNTRIES_URL' is not properly setup",
		})
		return
	}

	// Use Country Code to Call External API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("No Results found with country code - %s", countryCode),
		})
		return
	}

	// Map Country Data
	var externalData models.RestCountriesRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

// Weather and Such
func GetWeather(c *gin.Context) {
}

func GetTravelDestinations(c *gin.Context) {
}
