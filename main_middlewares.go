package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// SSOData ...
type SSOData struct {
	Timestamp string `json:"timestamp"`
	Token     string `json:"token"`
	AppSlug   string `json:"app_slug"`
}

func authenticateSharedToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authentication") != os.Getenv("SHARED_TOKEN") {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func authenticateWithSSO(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)

		var ssoData SSOData
		err := decoder.Decode(&ssoData)
		if err != nil {
			http.Error(w, "Failed to parse request", 500)
		}

		fmt.Println(ssoData.Timestamp)
		i, err := strconv.ParseInt(ssoData.Timestamp, 10, 64)
		if err != nil {
			panic(err)
		}
		requestTimestamp := time.Unix(i, 0)
		time := time.Now()
		fmt.Println(time)
		fmt.Println(requestTimestamp)
		if time.Sub(requestTimestamp).Seconds() > 30 {
			fmt.Println("outdated")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		s := sha1.New()
		s.Write([]byte(fmt.Sprintf("%s:%s:%s", ssoData.AppSlug, os.Getenv("SSO_TOKEN"), ssoData.Timestamp)))
		if hex.EncodeToString(s.Sum(nil)) != ssoData.Token {
			fmt.Println("token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
