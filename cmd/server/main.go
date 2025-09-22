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

	"zoubun/internal/middleware"
	"zoubun/internal/routes"

	_ "github.com/lib/pq"
)

func main() {
	// Database initialization
	db, err := sql.Open("postgres",
		fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable",
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_DB"),
		),
	)
	if err != nil {
		panic(err)
	}

	defer db.Close()
	services := routes.NewServices(db)
	//
	// Initialize the routes
	routes := routes.ConfigureRoutes(services)
	handler := middleware.ConfigureMiddleware(routes)
	log.Print("serving at localhost:3000")
	log.Panic(http.ListenAndServe(":3000", handler))
}
