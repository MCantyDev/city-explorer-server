package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

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
	Duration  time.Duration
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

	// Parse JWT_DURATION_DAYS
	JWTDuration, err := parseDurationWithDays(getEnv("JWT_DURATION"))
	if err != nil {
		log.Fatalf("invalid JWT_DURATION: %v", err)
	}
	if JWTDuration < time.Minute*15 || JWTDuration > (time.Hour*744) { // Ensure the Duration does not below 15 minutes or above 31 days
		log.Fatalf("invalid JWT_DURATION: Must be within range of 15m - 31d")
	}

	// Build Cfg
	Cfg.Database.User = getEnv("DB_USER")         // Initialise DB Server Username
	Cfg.Database.Password = getEnv("DB_PASSWORD") // Initialise DB Server Password
	Cfg.Database.Host = getEnv("DB_HOST")         // Initialise DB Server Host
	Cfg.Database.Port = port                      // Initialise DB Server Port (Parsed to Int)
	Cfg.Database.Name = getEnv("DB_NAME")         // Initialise DB Name

	Cfg.JWT.SecretKey = []byte(getEnv("JWT_SECRET_KEY")) // Initialise JWT Secret Key for JWT Encoding
	Cfg.JWT.Duration = JWTDuration                       // Initialise JWT Duration for Validation

	Cfg.PhotonAPI = ExternalAPI{
		Name: "Photon API",
		URL:  "https://photon.komoot.io/api/?q=%s&osm_tag=place:city&osm_tag=place:municipality&osm_tag=place:town&limit=50", // Static URL
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

// Parse Time with Days (As is the main type of time, however allow use of other suffixes such as 's', 'm', 'h', etc)
func parseDurationWithDays(duration string) (time.Duration, error) {
	if strings.HasSuffix(duration, "d") { // Custom check for Days suffix
		days := strings.TrimSuffix(duration, "d")

		daysInt, err := strconv.Atoi(days)
		if err != nil {
			log.Printf("Invalid day duration format: %v", err)
			return 0, fmt.Errorf("invalid days format: %v", err)
		}

		return time.ParseDuration(fmt.Sprintf("%dh", daysInt*24)) // Returns a (time.Duration, error)
	}
	return time.ParseDuration(duration)
}
