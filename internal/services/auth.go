package services

/* Authorisation Services

- Password Hasher
- JWT Token Validator
*/

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(unhashedPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(unhashedPassword), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CompareHashed(unhashedString string, hashedString string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(unhashedString))
	return err == nil
}
