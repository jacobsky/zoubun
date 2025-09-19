// Package router implements the generalized routing logic for the service
package router

import (
	"net/http"

	api "zoubun/internal/api"
)

func ConfigureRoutes() *http.ServeMux {
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", api.Index)
	routes.HandleFunc("GET /count", api.Count)
	routes.HandleFunc("POST /increment", api.Increment)
	routes.HandleFunc("POST /register", api.Register)
	routes.HandleFunc("POST /verify", api.Verify)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	return routes
}
