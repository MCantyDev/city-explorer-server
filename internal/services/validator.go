package services

import (
	"fmt"
	"unicode"

	"github.com/MCantyDev/city-explorer-server/internal/models"
)

var requirements models.UserCreationErrorRequirements = models.NewUserCreationErrorRequirements()

// Validate First name based on User Creation Error Requirements
func ValidateFirstName(firstName string) (string, bool) {
	if len(firstName) < requirements.MinFirstNameLength || len(firstName) > requirements.MaxFirstNameLength {
		return fmt.Sprintf("First name Must Be Between (%d - %d) Characters Long.",
			requirements.MinFirstNameLength, requirements.MaxFirstNameLength), false
	}
	return "", true
}

// Validate Last name based on User Creation Error Requirements
func ValidateLastName(lastName string) (string, bool) {
	if len(lastName) < requirements.MinLastNameLength || len(lastName) > requirements.MaxLastNameLength {
		return fmt.Sprintf("Last name Must Be Between (%d - %d) Characters Long.",
			requirements.MinLastNameLength, requirements.MaxLastNameLength), false
	}
	return "", true
}

// Validate Password based on User Creation Error Requirements
func ValidatePassword(password string) (string, bool) {
	if len(password) == 0 {
		return "Password Cannot Be Empty", false
	}

	if len(password) < requirements.MinPasswordLength || len(password) > requirements.MaxPasswordLength {
		return fmt.Sprintf("Password Must Be Between (%d - %d) Characters Long.", requirements.MinPasswordLength, requirements.MaxPasswordLength), false
	}

	if requirements.PasswordRequiresNumbers {
		if isValid := checkForNumbers(password); !isValid {
			return "Password Must Contain Numbers", false
		}
	}

	if requirements.PasswordRequiresSymbols {
		if isValid := checkForSymbols(password); !isValid {
			return "Password Must Contain Symbols", false
		}
	}
	return "", true
}

func checkForSymbols(word string) bool {
	for _, letter := range word {
		if !unicode.IsLetter(letter) && !unicode.IsDigit(letter) {
			return true
		}
	}
	return false
}

func checkForNumbers(word string) bool {
	for _, letter := range word {
		if unicode.IsDigit(letter) {
			return true
		}
	}
	return false
}
