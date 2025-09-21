package main

import (
	//	"math/big" // Used for big int counting for when this is hosted to be silly
	//	Used to track user contribution towards counting in a sqlite database.
	//	"gorm.io/driver/sqlite"
	//	"gorm.io/gorm"

	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	sqlc "zoubun/internal/db"
	"zoubun/internal/middleware"
	"zoubun/internal/routes"
)

func main() {
	// initialize the database here once the SQL is determined (should be easy)
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%v@%v@%v/%v", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_USER")))
	if err != nil {
		panic(err)
	}
	services := routes.NewServices(sqlc.New(db))
	// Initialize the routes
	routes := routes.ConfigureRoutes(services)
	handler := middleware.ConfigureMiddleware(routes)
	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", handler))
}
