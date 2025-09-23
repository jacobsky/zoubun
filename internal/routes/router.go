// Package routes implements the generalized routing logic for the service
package routes

import (
	"net/http"
)

func ConfigureRoutes(s *Services) *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", s.MessageOfTheDay)
	routes.HandleFunc("GET /motd", s.MessageOfTheDay)
	routes.Handle("GET /count", s.Authorize(http.HandlerFunc(s.Count)))
	routes.Handle("PUT /increment", s.Authorize(http.HandlerFunc(s.Increment)))
	routes.HandleFunc("POST /register", s.Register)
	routes.Handle("POST /rotate_key", s.Authorize(http.HandlerFunc(s.RotateKey)))
	routes.HandleFunc("GET /healthcheck", s.HealthCheck)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
