package utils_test

import (
	"pvz/internal/utils"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndGetRole(t *testing.T) {
	tests := []struct {
		name         string
		role         string
		modifyToken  func(string) string
		shouldError  bool
		expectedRole string
	}{
		{
			name:         "Valid token",
			role:         "admin",
			modifyToken:  func(token string) string { return token },
			shouldError:  false,
			expectedRole: "admin",
		},
		{
			name:        "Invalid token string",
			role:        "user",
			modifyToken: func(string) string { return "not.a.valid.token" },
			shouldError: true,
		},
		{
			name: "Expired token",
			role: "manager",
			modifyToken: func(_ string) string {
				// Генерируем вручную токен с истекшей датой
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
					"role":        "manager",
					"expire_date": time.Now().Add(-time.Hour).Unix(),
				})
				signed, _ := token.SignedString([]byte(utils.JwtSecret))
				return signed
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := utils.GenerateToken(tt.role)
			if err != nil {
				t.Fatalf("unexpected error during token generation: %v", err)
			}

			modifiedToken := tt.modifyToken(token)

			role, err := utils.GetRole(modifiedToken)
			if tt.shouldError {
				if err == nil {
					t.Errorf("expected error, got none")
				}
			} else {
				if err != nil {
					t.Errorf("did not expect error, got: %v", err)
				}
				if role != tt.expectedRole {
					t.Errorf("expected role %q, got %q", tt.expectedRole, role)
				}
			}
		})
	}
}
