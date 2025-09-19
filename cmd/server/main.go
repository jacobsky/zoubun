package main

import (
	//	"math/big" // Used for big int counting for when this is hosted to be silly
	//	Used to track user contribution towards counting in a sqlite database.
	//	"gorm.io/driver/sqlite"
	//	"gorm.io/gorm"

	"log"
	"net/http"

	"zoubun/internal/middleware"
	"zoubun/internal/router"
)

func main() {
	// initialize the database here once the SQL is determined (should be easy)
	// db := ?
	// Initialize the routes
	routes := router.ConfigureRoutes()
	handler := middleware.ConfigureMiddleware(routes)
	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", handler))
}
