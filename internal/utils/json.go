package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

func WriteJSONResponse(w http.ResponseWriter, status int, res any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(res)
}

func WriteJSONErrorResponse(w http.ResponseWriter, status int, message string) {
	WriteJSONResponse(w, status, map[string]string{
		"error": message,
	})
}

func DecodeJSON(r io.Reader, v any) error {
	decoder := json.NewDecoder(r)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}
