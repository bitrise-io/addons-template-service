package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter().StrictSlash(true)

	provision := r.PathPrefix("/provision").Subrouter()

	// Authenticate with the token set in env SHARED_TOKEN
	provision.Use(authenticateSharedToken)

	// POST /provision
	provision.HandleFunc("", postProvision).Methods(http.MethodPost)
	// PUT /provision/{app_slug}
	provision.HandleFunc("/{app_slug}", putProvision).Methods(http.MethodPut)
	// DELETE /provision/{app_slug}
	provision.HandleFunc("/{app_slug}", deleteProvision).Methods(http.MethodDelete)

	// will be available on localhost:5000
	http.ListenAndServe(":5000", r)
}
