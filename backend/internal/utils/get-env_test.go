package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"pvz/internal/utils"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name          string
		envKey        string
		envValue      string
		defaultValue  string
		expectedValue string
	}{
		{
			name:          "Env variable is set",
			envKey:        "TEST_KEY",
			envValue:      "test_value",
			defaultValue:  "default_value",
			expectedValue: "test_value",
		},
		{
			name:          "Env variable is not set",
			envKey:        "NON_EXISTENT_KEY",
			envValue:      "",
			defaultValue:  "default_value",
			expectedValue: "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			result := utils.GetEnv(tt.envKey, tt.defaultValue)
			assert.Equal(t, tt.expectedValue, result, "Unexpected value for "+tt.envKey)
		})
	}
}
