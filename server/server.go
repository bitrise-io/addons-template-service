package server

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	router     *mux.Router
	httpServer *http.Server
}

// Settings ...
type Settings struct {
	Port string
}

// NewServer ...
func NewServer(cfg Settings) *Server {
	cfg.fillWithDefaults()
	s := Server{
		router: mux.NewRouter().StrictSlash(false),
	}

	s.router.Path("/").Methods("GET").HandlerFunc(HomeHandler)

	prov := s.router.PathPrefix("/provision").Subrouter()
	prov.Use(accessTokenAuthMiddleware)
	prov.Methods("POST").HandlerFunc(PostProvisionHandler)
	prov.Path("/{app_slug:[0-9a-f]+}").Methods("DELETE").HandlerFunc(DeleteProvisionHandler)

	s.httpServer = &http.Server{
		Handler: s.router,
		Addr:    "0.0.0.0:" + cfg.Port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	s.router.NotFoundHandler = http.HandlerFunc(NotFoundHandler)

	return &s
}

func accessTokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.Debug("accessTokenAuthMiddleware")
		if strings.TrimPrefix(r.Header.Get("Authentication"), "token ") != os.Getenv("BITRISE_SHARED_SECRET") {
			if err := renderErrorMessage(w, http.StatusUnauthorized, "Unauthorized"); err != nil {
				fmt.Printf("failed to render error JSON, error: %s\n", err)
			}
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Start ...
func (s Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (cfg *Settings) fillWithDefaults() {
	if cfg.Port == "" {
		cfg.Port = "3000"
	}
}

// NotFoundHandler ...
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	fields := logrus.Fields{
		"method": r.Method,
		"path":   r.URL.Path,
		"status": http.StatusNotFound,
		"header": r.Header,
	}
	logrus.WithFields(fields).Warn("Not Found")
	if err := renderErrorMessage(w, http.StatusNotFound, "Not Found"); err != nil {
		fmt.Printf("failed to render error JSON, error: %s\n", err)
	}
}
