package utils

import (
	"fmt"
	"regexp"
)

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
