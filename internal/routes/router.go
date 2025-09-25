// Package routes implements the generalized routing logic for the service
package routes

import (
	"net/http"
)

func ConfigureRoutes(s *Services) *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", Handler(s.MessageOfTheDay))
	routes.HandleFunc("GET /motd", Handler(s.MessageOfTheDay))
	routes.Handle("GET /count", s.Authorize(Handler(s.Count)))
	routes.Handle("PUT /increment", s.Authorize(Handler(s.Increment)))
	routes.HandleFunc("POST /register", Handler(s.Register))
	routes.Handle("POST /rotate_key", s.Authorize(Handler(s.RotateKey)))
	routes.HandleFunc("GET /healthcheck", Handler(s.HealthCheck))

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
