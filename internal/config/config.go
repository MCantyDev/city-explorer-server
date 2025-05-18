package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Cfg *Config

type Config struct {
	// DB Data
	Database DatabaseConfig

	// JWT Data
	JWT JWTConfig

	// External API URLs
	PhotonAPI        ExternalAPI
	RestCountriesAPI ExternalAPI
	OpenWeatherAPI   SecureExternalAPI
	OpenTripAPI      SecureExternalAPI
	OpenTripXIDAPI   SecureExternalAPI
}

type DatabaseConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	Name     string
}

type JWTConfig struct {
	SecretKey []byte
}

type ExternalAPI struct {
	Name string
	URL  string
}

type SecureExternalAPI struct {
	Name string
	URL  string
	Key  string
}

// Load critical config values and initialise them
func Load() {
	Cfg = &Config{} // Initialise Cfg before Assigning any values (or i get a runtime error)

	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found, use the .env.example to build a .env file in the root directory")
	}

	// Parse DB_PORT
	port, err := strconv.Atoi(getEnv("DB_PORT"))
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	if port < 1024 || port > 65535 {
		log.Fatal("invalid DB_PORT: Must be within range of 1024 - 65535")
	}

	// Build Cfg
	Cfg.Database.User = getEnv("DB_USER")         // Initialise DB Server Username
	Cfg.Database.Password = getEnv("DB_PASSWORD") // Initialise DB Server Password
	Cfg.Database.Host = getEnv("DB_HOST")         // Initialise DB Server Host
	Cfg.Database.Port = port                      // Initialise DB Server Port (Parsed to Int)
	Cfg.Database.Name = getEnv("DB_NAME")         // Initialise DB Name

	Cfg.JWT.SecretKey = []byte(getEnv("JWT_SECRET_KEY")) // Initialise JWT Secret Key for JWT Encoding

	Cfg.PhotonAPI = ExternalAPI{
		Name: "Photon API",
		URL:  "https://photon.komoot.io/api/?q=%s&lang=en", // Static URL
	}
	Cfg.RestCountriesAPI = ExternalAPI{
		Name: "Rest-Countries API",
		URL:  "https://restcountries.com/v3.1/alpha/%s", // Static URL
	}
	Cfg.OpenWeatherAPI = SecureExternalAPI{
		Name: "OpenWeather API",
		URL:  "https://api.openweathermap.org/data/3.0/onecall?lat=%s&lon=%s&exclude=alerts,hourly,minutely&units=metric&appid=%s", // Static URL
		Key:  getEnv("OPENWEATHER_KEY"),
	}
	Cfg.OpenTripAPI = SecureExternalAPI{
		Name: "OpenTrip API",
		URL:  "https://api.opentripmap.com/0.1/en/places/radius?lat=%s&lon=%s&radius=2000&limit=50&kinds=amusements,accomodations,tourist_facilities&rate=2&apikey=%s", // Static URL
		Key:  getEnv("OPENTRIP_KEY"),
	}
	Cfg.OpenTripXIDAPI = SecureExternalAPI{
		Name: "OpenTrip Singular Place API",
		URL:  "https://api.opentripmap.com/0.1/en/places/xid/%s?apikey=%s",
		Key:  getEnv("OPENTRIP_KEY"),
	}
}

// getEnv returns the environment variable or shuts down if missing

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required environment variable: %s", key) // force server shutdown if env variable is missing
	}
	return val
}
