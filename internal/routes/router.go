// Package routes implements the generalized routing logic for the service
package routes

import (
	"net/http"

	"zoubun/internal/middleware"
)

func ConfigureRoutes(s *Services) *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", s.MessageOfTheDay)
	routes.HandleFunc("GET /motd", s.MessageOfTheDay)
	routes.Handle("GET /count", middleware.Authorize(http.HandlerFunc(s.Count)))
	routes.Handle("POST /regenerate_key", middleware.Authorize(http.HandlerFunc(s.Count)))
	routes.Handle("PUT /increment", middleware.Authorize(http.HandlerFunc(s.Increment)))
	routes.HandleFunc("POST /register", s.Register)
	routes.HandleFunc("GET /healthcheck", s.HealthCheck)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
