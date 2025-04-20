package utils

func ValidateRole(role string) bool {
	if role != "moderator" && role != "client" {
		return false
	}

	return true
}
