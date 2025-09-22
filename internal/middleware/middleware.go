// Package middleware implements the middleware used by the server
package middleware

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AuthorizationError struct {
	Kind    string `json:"kind"`
	Details string `json:"details"`
}

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// TODO: Placeholder for HTTP authorization middleware
		// Check the key against the DB
		// If the key is valid, then serve the request
		key := req.Header.Get("zoubun-api-key")
		if key == "" {
			resp.WriteHeader(401)
			json.NewEncoder(resp).Encode(AuthorizationError{
				Kind:    "Unauthorized",
				Details: "No authorization header was found. Please ensure that the API key is included in the `zoubun-api-key` header",
			})
			// TODO: Put a better error message in here
			resp.Write(make([]byte, 0))
			return
		}
		// TODO: Plumb it into a DB query
		id := 1
		if key != "hirakegoma" {
			// TODO: Add in the hashing/SQL logic here
			resp.WriteHeader(403)
			// TODO: Put a better error message in here
			json.NewEncoder(resp).Encode(AuthorizationError{
				Kind:    "Forbidden",
				Details: "The `zoubun-api-key` header is invalid or expired. Please update your key",
			})
			resp.Write(make([]byte, 0))
			return
		}
		// Add the Userid to the request header so that it's easier to track
		req.Header.Add("userid", strconv.Itoa(id))
		next.ServeHTTP(resp, req)
		// Otherwise return no
	})
}

var httpRequestCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_total",
	Help: "Total number of HTTP requests received",
}, []string{"status", "path", "method"})

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rec *statusRecorder) WriteHEader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Middleware that provides some various request specific metrics to prometheus
func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		recorder := statusRecorder{
			ResponseWriter: resp,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, req)

		method := req.Method
		path := req.URL.Path
		status := strconv.Itoa(recorder.statusCode)

		httpRequestCounter.WithLabelValues(status, path, method).Inc()
	})
}

// Adds the request ID to each request
func addRequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// TODO: Add a response header
		reqid := uuid.New()
		resp.Header().Add("Zoubun-Request-ID", reqid.String())
		next.ServeHTTP(resp, req)
	})
}

func ConfigureMiddleware(routes *http.ServeMux) http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(httpRequestCounter)
	reqHandler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)

	routes.Handle("/metrics", reqHandler)

	handler := prometheusMiddleware(routes)
	handler = addRequestIDMiddleware(handler)
	return handler
}
