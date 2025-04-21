package utils_test

import (
	"pvz/internal/utils"
	"testing"
	"time"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid email", "test@example.com", true},
		{"Missing @", "testexample.com", false},
		{"Missing domain", "test@", false},
		{"Missing local part", "@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ValidateEmail(tt.email); got != tt.expected {
				t.Errorf("ValidateEmail(%q) = %v, want %v", tt.email, got, tt.expected)
			}
		})
	}
}

func TestValidateRole(t *testing.T) {
	tests := []struct {
		name     string
		role     string
		expected bool
	}{
		{"Valid moderator", "moderator", true},
		{"Valid client", "client", true},
		{"Invalid role", "admin", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utils.ValidateRole(tt.role); got != tt.expected {
				t.Errorf("ValidateRole(%q) = %v, want %v", tt.role, got, tt.expected)
			}
		})
	}
}

func TestValidateCity(t *testing.T) {
	tests := []struct {
		name      string
		city      string
		expectErr bool
	}{
		{"Valid city", "Казань", false},
		{"Invalid city", "Томск", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateCity(tt.city)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateCity(%q) error = %v, wantErr %v", tt.city, err, tt.expectErr)
			}
		})
	}
}

func TestValidateProductType(t *testing.T) {
	tests := []struct {
		name      string
		ptype     string
		expectErr bool
	}{
		{"Valid type", "одежда", false},
		{"Invalid type", "мебель", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateProductType(tt.ptype)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateProductType(%q) error = %v, wantErr %v", tt.ptype, err, tt.expectErr)
			}
		})
	}
}

func TestValidateTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name      string
		start     time.Time
		end       time.Time
		expectErr bool
	}{
		{"Valid time", now, now.Add(1 * time.Hour), false},
		{"Start after end", now.Add(1 * time.Hour), now, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateTime(tt.start, tt.end)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateTime(%v, %v) error = %v, wantErr %v", tt.start, tt.end, err, tt.expectErr)
			}
		})
	}
}

func TestValidateAll(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		role      string
		expectErr bool
	}{
		{"Valid all", "test@test.com", "client", false},
		{"Invalid email", "invalid", "client", true},
		{"Invalid role", "test@test.com", "admin", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := utils.ValidateAll(tt.email, tt.role)
			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateAll(%q, %q) error = %v, wantErr %v", tt.email, tt.role, err, tt.expectErr)
			}
		})
	}
}
