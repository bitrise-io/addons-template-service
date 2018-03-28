package main

import (
	"net/http"
	"os"
)

func authenticateSharedToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authentication") != os.Getenv("SHARED_TOKEN") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
