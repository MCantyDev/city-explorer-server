package services

import (
	"github.com/MCantyDev/city-explorer-server/internal/database"
	"github.com/MCantyDev/city-explorer-server/internal/models"
)

func GetOrCreateCity(name string) (*models.City, error) {
	var city models.City

	query := database.NewQueryBuilder("SELECT").Table("cities").Where("name = ?").Build()
	_, err := database.Execute(&city, query, name)
	if err == nil && city.Id > 0 {
		return &city, nil
	}

	city = models.City{Name: name}
	_, err = database.Execute(&city, "INSERT")
	if err != nil {
		return nil, err
	}

	return &city, nil
}

func GetCountry(code string) (*models.Country, error) {
	var country models.Country

	query := database.NewQueryBuilder("SELECT").Table("countries").Where("iso_code = ?").Build()
	_, err := database.Execute(&country, query, code)
	if err == nil && country.Id > 0 {
		return &country, nil
	}

	return nil, err
}

// Interfaces are weird...so it made userId a float64...so i rolled with it...
func CheckAdminStatus(userId float64) (bool, error) {
	var user models.User

	query := database.NewQueryBuilder("SELECT").Table("users").Where("id = ?").Build()
	_, err := database.Execute(&user, query, userId)
	if err == nil && user.Id > 0 {
		return user.IsAdmin, nil
	}

	return false, err
}
