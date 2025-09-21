// Package routes implements the generalized routing logic for the service
package routes

import (
	"net/http"
)

func ConfigureRoutes(s *Services) *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", s.Index)
	routes.HandleFunc("GET /count", s.Count)
	routes.HandleFunc("POST /increment", s.Increment)
	routes.HandleFunc("POST /register", s.Register)
	routes.HandleFunc("POST /verify", s.Verify)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
