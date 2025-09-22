// Package routes implements the generalized routing logic for the service
package routes

import (
	"net/http"
)

func ConfigureRoutes(s *Services) *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", s.MessageOfTheDay)
	routes.HandleFunc("GET /motd", s.MessageOfTheDay)
	routes.HandleFunc("GET /{userid}/count/", s.Count)
	routes.HandleFunc("PUT /{userid}/increment", s.Increment)
	routes.HandleFunc("POST /register", s.Register)
	routes.HandleFunc("GET /healthcheck", s.HealthCheck)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
