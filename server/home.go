package server

import (
	"fmt"
	"net/http"
	"time"
)

// HomeHandler ...
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeResponse := map[string]interface{}{}
	homeResponse["time"] = time.Now().UTC().Format(time.RFC3339)
	homeResponse["message"] = "Bitrise Test Results Addon Server"

	if err := renderJSON(w, http.StatusOK, homeResponse); err != nil {
		fmt.Printf("failed to render JSON, error: %s\n", err)
	}
}
