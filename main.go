package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

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

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	log.Printf("Server started at port %v", port)
	err = http.ListenAndServe(":"+port, r)
	log.Fatal(err)
}
