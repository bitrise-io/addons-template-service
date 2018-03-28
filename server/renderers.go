package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func renderJSON(w http.ResponseWriter, status int, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return renderErrorMessage(w, http.StatusInternalServerError, err.Error())
	}

	w.WriteHeader(status)
	if _, err := w.Write(jsonResp); err != nil {
		fmt.Printf("failed to write response, error: %s\n", err)
	}
	return nil
}

func renderErrorMessage(w http.ResponseWriter, status int, errorMessage string) error {
	return renderJSON(w, status, map[string]string{"error": errorMessage})
}
