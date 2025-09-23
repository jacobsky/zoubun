// Package middleware implements the middleware used by the server
package middleware

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

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
