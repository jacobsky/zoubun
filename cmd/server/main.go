package main

import (
	//	"math/big" // Used for big int counting for when this is hosted to be silly
	//	Used to track user contribution towards counting in a sqlite database.
	//	"gorm.io/driver/sqlite"
	//	"gorm.io/gorm"

	"encoding/json"
	"log"
	"net/http"
)

type Counter struct {
	Count int `json:"count"`
}

var currentCount = Counter{0}

func index(resp http.ResponseWriter, req *http.Request) {
	const Endpoint = "/incrementや/countで増分してみよう～"
	resp.Header().Set("ContentType", "text/html; charset=utf-8")
	resp.Write([]byte(Endpoint))
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

func main() {
	// initialize the database here once the SQL is determined (should be easy)
	// db := ?
	// Initialize the routes
	routes := http.NewServeMux()
	routes.HandleFunc("GET /", index)
	routes.HandleFunc("GET /count", count)
	routes.HandleFunc("POST /increment", increment)

	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", routes))
}
