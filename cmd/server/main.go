package main

import (
	//	"math/big" // Used for big int counting for when this is hosted to be silly
	//	Used to track user contribution towards counting in a sqlite database.
	//	"gorm.io/driver/sqlite"
	//	"gorm.io/gorm"
	"encoding/json"
	"log"
	"net/http"

	models "zoubun/internal"
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

func main() {
	// initialize the database here once the SQL is determined (should be easy)
	// db := ?
	// Initialize the routes
	routes := http.NewServeMux()
	routes.HandleFunc("GET /", index)
	routes.HandleFunc("GET /count", count)
	routes.HandleFunc("POST /increment", increment)
	routes.HandleFunc("POST /register", register)
	routes.HandleFunc("POST /verify", verify)
	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", routes))
}
