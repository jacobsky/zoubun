// Package middleware implements the middleware used by the server
package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// TODO: Placeholder for HTTP authorization middleware
		// Check the key against the DB
		// If the key is valid, then serve the request
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

func ConfigureMiddleware(routes *http.ServeMux) http.Handler {
	reg := prometheus.NewRegistry()
	reg.MustRegister(httpRequestCounter)
	reqHandler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)

	routes.Handle("/metrics", reqHandler)

	handler := prometheusMiddleware(routes)
	return handler
}
