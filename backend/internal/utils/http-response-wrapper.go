package utils

import (
	"encoding/json"
	"net/http"
	
	"pvz/internal/delivery/forms"
)

func WriteJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(forms.ErrorForm{Message: message})
}

func WriteJson(w http.ResponseWriter, content interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(content)
}
