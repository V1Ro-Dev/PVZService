package utils

import (
	"errors"
	"fmt"
	"regexp"
)

var allowedCities = []string{"Москва", "Санкт-Петербург", "Казань"}

func ValidateRole(role string) bool {
	if role != "moderator" && role != "client" {
		return false
	}

	return true
}

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._-]+@[a-zA-Z0-9._-]+\.[a-zA-Z0-9_-]+$`

	re := regexp.MustCompile(regex)

	return re.MatchString(email)
}

func ValidateAll(email string, role string) error {
	if !ValidateEmail(email) {
		return fmt.Errorf("email %s is not valid", email)
	}
	if !ValidateRole(role) {
		return fmt.Errorf("role %s is not valid", role)
	}

	return nil
}

func ValidateCity(city string) error {
	for _, valid := range allowedCities {
		if city == valid {
			return nil
		}
	}

	return errors.New("invalid city")
}
