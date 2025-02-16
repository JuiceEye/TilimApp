package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func ParseJSONRequest(r *http.Request, body any) error {
	if r.Body == nil {
		return fmt.Errorf("missing body")
	}
	return json.NewDecoder(r.Body).Decode(body)
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	log.Printf("Response: %s", payload)
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	_ = WriteJSONResponse(w, status, map[string]string{"error": err.Error()})
}
