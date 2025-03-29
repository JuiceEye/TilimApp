package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"tilimauth/internal/dto/request"
)

func ParseJSONRequest(r *http.Request, body any) error {
	if r.ContentLength == 0 {
		return fmt.Errorf("missing body")
	}

	err := json.NewDecoder(r.Body).Decode(body)

	jsonBody, _ := json.MarshalIndent(body, "", "  ")
	log.Printf("Parsing request %v", string(jsonBody))

	return err
}

func ParseAndValidate(r *http.Request, req request.Request) error {
	if err := ParseJSONRequest(r, req); err != nil {
		return errors.New("invalid JSON request")
	}

	if err := req.Validate(); err != nil {
		return err
	}

	return nil
}

func WriteJSONResponse(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	jsonPayload, _ := json.MarshalIndent(payload, "", "  ")
	log.Printf("Response: %v \n", string(jsonPayload))
	fmt.Printf("-----------------------------------------------------------------------------------------------")
	return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	_ = WriteJSONResponse(w, status, map[string]string{"error": err.Error()})
}
