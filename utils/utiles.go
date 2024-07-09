package utils

import (
	"encoding/json"
	"net/http"

	"gopkg.in/go-playground/validator.v9"
)

type Validator interface {
	Struct(s interface{}) error
}

var Validate Validator = validator.New()

func WriteJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ReadJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func WriteError(w http.ResponseWriter, statusCode int, message string) {
	WriteJSON(w, statusCode, map[string]string{"error": message})
}
