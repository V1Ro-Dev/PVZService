package utils

import (
	"errors"
	"fmt"
	"pvz/internal/models"
	"regexp"
)

var allowedCities = []string{"Москва", "Санкт-Петербург", "Казань"}
var allowedTypes = []string{"электроника", "одежда", "обувь"}

func ValidateRole(role string) bool {
	if role != string(models.Moderator) && role != string(models.Client) && role != string(models.Employee) {
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

func ValidateProductType(productType string) error {
	for _, valid := range allowedTypes {
		if productType == valid {
			return nil
		}
	}

	return errors.New("invalid product type")
}
