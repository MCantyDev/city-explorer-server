package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/MCantyDev/city-explorer-server/internal/config"
	"github.com/MCantyDev/city-explorer-server/internal/database"
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
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

func GetCountry(c *gin.Context) {
	countryCode := c.Query("country-code")
	name := c.Query("name")

	// Ensure Country Code was set in the Get Request
	if countryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'country-code' query parameter",
		})
		return
	}
	// Ensure Country Code is of ISO 3166-1 Alpha-2 type
	if len(countryCode) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Must submit ISO 3166-1 Alpha-2 country code",
		})
		return
	}

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'name' query parameter",
		})
		return
	}

	// Try to Retrieve Data from Server (IF nothing is returned retrieve using External API)
	var countryData models.Country
	query := database.NewQueryBuilder("SELECT").Table("countries").Where("iso_code = ?").Build()
	_, err := database.Execute(&countryData, query, countryCode)
	if err == nil && countryData.Id > 0 && countryData.ExpiryDate.After(time.Now()) {
		c.JSON(http.StatusOK, countryData.Data)
		return
	}

	// Get the URL Template for Rest Countries API
	url := fmt.Sprintf(config.Cfg.RestCountriesAPI.URL, countryCode)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.RestCountriesAPI.Name + "' URL is not properly setup",
		})
		return
	}

	// Use Country Code to Call External API
	var externalData models.RestCountriesRequest
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Map Country Data
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	expiry := time.Now().AddDate(1, 0, 0)

	if countryData.Id > 0 && countryData.ExpiryDate.Before(time.Now()) {
		query = database.NewQueryBuilder("UPDATE").Table("countries").Columns("data", "expiry_date").Where("id = ?").Build()
		_, err = database.Execute(nil, query, externalData, expiry, countryData.Id)
	} else {
		query = database.NewQueryBuilder("INSERT").Table("countries").Columns("name", "iso_code", "data", "expiry_date").Values(4).Build()
		_, err = database.Execute(nil, query, name, countryCode, externalData, time.Now().AddDate(1, 0, 0))
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while saving weather data",
		})
		return
	}
	c.JSON(http.StatusOK, externalData[0])
}

func GetWeather(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")
	if lat == "" || long == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' or 'long' query parameter",
		})
		return
	}

	city := c.Query("city")
	country := c.Query("country-code")
	if city == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'city' or 'country' query parameter",
		})
	}

	// Try to retrieve data from the database
	var cityWeather models.CityWeather
	query := database.NewQueryBuilder("SELECT").Table("city_weather").Where("ABS(lat - ?) < 0.0001 AND ABS(lon - ?) < 0.0001").Build()
	_, err := database.Execute(&cityWeather, query, lat, long)
	if err == nil && cityWeather.Id > 0 && cityWeather.ExpiryDate.After(time.Now()) {
		c.JSON(http.StatusOK, cityWeather.Data)
		return
	}

	// Build URL String for OpenWeatherAPI
	url := fmt.Sprintf(config.Cfg.OpenWeatherAPI.URL, lat, long, config.Cfg.OpenWeatherAPI.Key)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenWeatherAPI.Name + "' URL is not set up properly",
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

	// Map City Weather Data
	var externalData models.OpenWeatherRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	expiry := time.Now().AddDate(0, 0, 1)

	if cityWeather.Id > 0 {
		query = database.NewQueryBuilder("UPDATE").Table("city_weather").Columns("data", "expiry_date").Where("id = ?").Build()
		_, err = database.Execute(nil, query, data, expiry, cityWeather.Id)
	} else {
		city, _ := services.GetOrCreateCity(city)
		country, _ := services.GetCountry(country)
		query = database.NewQueryBuilder("INSERT").Table("city_weather").Columns("lat", "lon", "city_id", "country_id", "data", "expiry_date").Values(6).Build()
		_, err = database.Execute(nil, query, lat, long, city.Id, country.Id, data, expiry)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while saving weather data",
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

func GetTravelDestinations(c *gin.Context) {
	lat := c.Query("lat")
	long := c.Query("long")
	if lat == "" || long == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'lat' or 'long' query parameter",
		})
		return
	}

	city := c.Query("city")
	countryCode := c.Query("country-code")
	if city == "" || countryCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'city' or 'country' query parameter",
		})
		return
	}

	// Try to retrieve from DB
	var citySights models.CitySights
	query := database.NewQueryBuilder("SELECT").Table("city_sights").Where("ABS(lat - ?) < 0.0001 AND ABS(lon - ?) < 0.0001").Build()

	_, err := database.Execute(&citySights, query, lat, long)
	if err == nil && citySights.Id > 0 && !citySights.ExpiryDate.IsZero() && citySights.ExpiryDate.After(time.Now()) {
		c.JSON(http.StatusOK, citySights.Data)
		return
	}

	// Build URL String for OpenTripAPI
	url := fmt.Sprintf(config.Cfg.OpenTripAPI.URL, lat, long, config.Cfg.OpenTripAPI.Key)
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

	// Map OpenTripRequest Data
	var externalData models.OpenTripRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	expiry := time.Now().AddDate(0, 6, 0)

	// Either update or insert new row
	if citySights.Id > 0 {
		query = database.NewQueryBuilder("UPDATE").Table("city_sights").Columns("data", "expiry_date").Where("id = ?").Build()
		_, err = database.Execute(nil, query, data, expiry, citySights.Id)
	} else {
		city, _ := services.GetOrCreateCity(city)
		country, _ := services.GetCountry(countryCode)
		query = database.NewQueryBuilder("INSERT").Table("city_sights").Columns("lat", "lon", "city_id", "country_id", "data", "expiry_date").Values(6).Build()
		_, err = database.Execute(nil, query, lat, long, city.Id, country.Id, data, expiry)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while saving travel destinations",
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}

func GetTravelDestination(c *gin.Context) {
	xid := c.Query("xid")
	if xid == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'xid' query parameter",
		})
		return
	}

	city := c.Query("city")
	country := c.Query("country-code")
	if city == "" || country == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Missing 'city' or 'country' query parameter",
		})
	}

	// Try to retrieve from DB first
	var poiData models.CityPoi
	query := database.NewQueryBuilder("SELECT").Table("city_pois").Where("xid = ?").Build()
	_, err := database.Execute(&poiData, query, xid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error (Try again later)",
		})
		return
	}

	// If data found and not expired, return it
	if poiData.Id > 0 && poiData.ExpiryDate.After(time.Now()) {
		c.JSON(http.StatusOK, poiData.Data)
		return
	}

	// Build API URL
	url := fmt.Sprintf(config.Cfg.OpenTripXIDAPI.URL, xid, config.Cfg.OpenTripAPI.Key)
	if url == "" || !strings.HasPrefix(url, "https://") {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "'" + config.Cfg.OpenTripXIDAPI.Name + "' URL is not setup properly",
		})
		return
	}

	// Fetch from API
	data, err := services.FetchExternalAPI(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Unmarshal just to validate
	var externalData models.OpenTripPlaceRequest
	err = json.Unmarshal(data, &externalData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	expiry := time.Now().AddDate(0, 6, 0)

	if poiData.Id > 0 {
		query = database.NewQueryBuilder("UPDATE").Table("city_pois").Columns("data", "expiry_date").Where("id = ?").Build()
		_, err = database.Execute(nil, query, data, expiry, poiData.Id)
	} else {
		city, _ := services.GetOrCreateCity(city)
		country, _ := services.GetCountry(country)
		query = database.NewQueryBuilder("INSERT").Table("city_pois").Columns("city_id", "country_id", "xid", "data", "expiry_date").Values(5).Build()
		_, err = database.Execute(nil, query, city.Id, country.Id, xid, data, expiry)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Database error while saving travel Pois",
		})
		return
	}

	c.JSON(http.StatusOK, externalData)
}
