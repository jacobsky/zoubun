package main

import (
	//	"math/big" // Used for big int counting for when this is hosted to be silly
	//	Used to track user contribution towards counting in a sqlite database.
	//	"gorm.io/driver/sqlite"
	//	"gorm.io/gorm"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	models "zoubun/internal"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var currentCount = models.Counter{Count: 0}

func index(resp http.ResponseWriter, req *http.Request) {
	motd := models.Motd{Message: "皆さん、/incrementや/countのエンドポイントで増分してみよう～"}
	resp.Header().Set("ContentType", "text/html; charset=utf-8")
	json.NewEncoder(resp).Encode(motd)
}

func count(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func increment(resp http.ResponseWriter, req *http.Request) {
	currentCount.Count++
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func register(resp http.ResponseWriter, req *http.Request) {
	// TODO: This function will be used to register a user with a unique key.
}

func verify(resp http.ResponseWriter, req *http.Request) {
	// TODO: The verification endpoint that
}

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

func main() {
	// initialize the database here once the SQL is determined (should be easy)
	// db := ?
	// Initialize the routes
	routes := http.NewServeMux()
	routes.HandleFunc("GET /index", index)
	routes.HandleFunc("GET /count", count)
	routes.HandleFunc("POST /increment", increment)
	routes.HandleFunc("POST /register", register)
	routes.HandleFunc("POST /verify", verify)

	// This is specifically to handle all the emission of prometheus metrics for monitoring.
	reg := prometheus.NewRegistry()
	reg.MustRegister(httpRequestCounter)
	handler := promhttp.HandlerFor(
		reg, promhttp.HandlerOpts{},
	)
	routes.Handle("/metrics", handler)

	promHandler := prometheusMiddleware(routes)
	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", promHandler))
}
