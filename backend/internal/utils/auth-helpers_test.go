package utils_test

import (
	"testing"

	"pvz/internal/utils"
)

func TestHashPassword(t *testing.T) {
	// Таблица тестов для HashPassword
	tests := []struct {
		password string
		salt     string
		expected string
	}{
		{
			password: "testPassword",
			salt:     "randomSalt",
			expected: "7ca992378e2614b944c140ebfb9fe3b74b77ecaae14ff011f3dc5caa541566dc",
		},
		{
			password: "anotherPassword",
			salt:     "differentSalt",
			expected: "37d5a2fc25be7151bfc62bb809deb8ab197afcd1f884b230ffc3c49e4fb17935",
		},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			hashedPassword := utils.HashPassword(tt.password, tt.salt)
			if hashedPassword != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, hashedPassword)
			}
		})
	}
}

func TestGenSalt(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"test_salt_1"},
		{"test_salt_2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt := utils.GenSalt()

			if len(salt) != 10 {
				t.Errorf("expected salt of length 10, got %d", len(salt))
			}

			for _, char := range salt {
				if !isValidSaltChar(char) {
					t.Errorf("invalid character '%c' in salt", char)
				}
			}
		})
	}
}

func isValidSaltChar(c rune) bool {
	const validChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	for _, validChar := range validChars {
		if c == validChar {
			return true
		}
	}
	return false
}

func TestCheckPassword(t *testing.T) {
	tests := []struct {
		password     string
		userPassword string
		userSalt     string
		expected     bool
	}{
		{
			password:     "testPassword",
			userPassword: utils.HashPassword("testPassword", "randomSalt"),
			userSalt:     "randomSalt",
			expected:     true,
		},
		{
			password:     "wrongPassword",
			userPassword: utils.HashPassword("testPassword", "randomSalt"),
			userSalt:     "randomSalt",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.password, func(t *testing.T) {
			result := utils.CheckPassword(tt.password, tt.userPassword, tt.userSalt)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
