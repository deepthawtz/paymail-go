package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// ServiceDiscovery is returned when well known file is requested by a Paymail client
type ServiceDiscovery struct {
	Version      string                 `json:"bsvalias"`
	Capabilities map[string]interface{} `json:"capabilities"`
}

// Server is a Paymail server
type Server struct {
	BaseURL string
	Port    string
}

// NewServer creates a new Paymail server
func NewServer(baseURL, port string) *Server {
	return &Server{BaseURL: baseURL, Port: port}
}

// Start starts the Paymail server
func (s *Server) Start() {
	r := mux.NewRouter()
	r.Use(setContentTypeHeader)
	r.HandleFunc("/.well-known/bsvalias", loggingMiddleware(s.ServiceDiscoveryHandler)).Methods("GET")
	r.HandleFunc("/api/v1/bsvalias/id/{paymail}", loggingMiddleware(s.IdentityHandler)).Methods("GET")
	r.HandleFunc("/api/v1/bsvalias/address/{paymail}", loggingMiddleware(s.PaymentDestinationHandler)).Methods("GET")

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", s.Port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	logrus.Infof("Listening on port %s", s.Port)
	logrus.Fatal(srv.ListenAndServe())
}

// ServiceDiscoveryHandler handles request for Paymail server capabilities
func (s *Server) ServiceDiscoveryHandler(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(&ServiceDiscovery{
		Version: "1.0",
		Capabilities: map[string]interface{}{
			"pki":                s.BaseURL + "/api/v1/bsvalias/id/{alias}@{domain.tld}",
			"paymentDestination": s.BaseURL + "/api/v1/bsvalias/address/{alias}@{domain.tld}",
			// TODO: add capabilities
		},
	})

	if err != nil {
		logrus.Warn(err)
	}
	io.WriteString(w, string(j))
}

// IdentityHandler returns identity for a given paymail {alias}@{domain}.{tld}
func (s *Server) IdentityHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

// PaymentDestinationHandler returns payment destination for a given paymail {alias}@{domain}.{tld}
func (s *Server) PaymentDestinationHandler(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func setContentTypeHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "paymail-go")
		w.Header().Set("Content-Type", `application/json; schema="https://schemas.nchain.com/bsvalias/1.0/capability-discovery"`)
		next.ServeHTTP(w, r)
	})
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := time.Now()
		next.ServeHTTP(w, r)
		logrus.WithFields(logrus.Fields{
			"url":      r.URL.String(),
			"duration": fmt.Sprintf("%s", time.Since(s)),
		}).Info("served request")
	})
}
