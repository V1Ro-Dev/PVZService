package utils_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"pvz/internal/delivery/forms"
	"pvz/internal/utils"
)

func TestWriteJsonError(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		statusCode int
	}{
		{"Error 400", "Bad Request", 400},
		{"Error 404", "Not Found", 404},
		{"Error 500", "Internal Server Error", 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			utils.WriteJsonError(rr, tt.message, tt.statusCode)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, status)
			}

			expected := forms.ErrorForm{Message: tt.message}
			var actual forms.ErrorForm
			if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			if actual.Message != expected.Message {
				t.Errorf("expected message %q, got %q", expected.Message, actual.Message)
			}
		})
	}
}

func TestWriteJson(t *testing.T) {
	tests := []struct {
		name       string
		content    interface{}
		statusCode int
	}{
		{"Success", map[string]string{"message": "success"}, 200},
		{"Created", map[string]string{"message": "resource created"}, 201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := httptest.NewRecorder()

			utils.WriteJson(rr, tt.content, tt.statusCode)

			if status := rr.Code; status != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, status)
			}

			var actual map[string]string
			if err := json.NewDecoder(rr.Body).Decode(&actual); err != nil {
				t.Fatalf("failed to decode response body: %v", err)
			}

			for key, val := range tt.content.(map[string]string) {
				if actual[key] != val {
					t.Errorf("expected key %q to have value %q, got %q", key, val, actual[key])
				}
			}
		})
	}
}
