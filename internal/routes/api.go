package routes

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	sqlc "zoubun/internal/db"
)

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type Counter struct {
	Count int `json:"count"`
}

type Motd struct {
	Message string `json:"message"`
}

type RegisterRequest struct {
	Username string `json:"username"`
}

type RegisterResponse struct {
	APIKey          string `json:"api-key"`
	VerificationKey string `json:"verificationkey"`
}

var currentCount = Counter{Count: 0}

type Services struct {
	db      *sql.DB
	queries *sqlc.Queries
}

func NewServices(db *sql.DB) *Services {
	return &Services{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (s *Services) MessageOfTheDay(resp http.ResponseWriter, req *http.Request) {
	message, err := s.queries.Motd(req.Context())
	if err != nil {
		log.Fatalf("Database Error %v", err.Error())
		resp.WriteHeader(503)
		errResp := ErrorResponse{
			Message: "An internal database error ha occurred. If this issue persists, please notify the administrators",
		}
		json.NewEncoder(resp).Encode(errResp)
		return
	}

	motd := Motd{Message: message}
	json.NewEncoder(resp).Encode(motd)
}

func (s *Services) Count(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func (s *Services) Increment(resp http.ResponseWriter, req *http.Request) {
	currentCount.Count++
	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func (s *Services) Register(resp http.ResponseWriter, req *http.Request) {
	// TODO: This function will be used to register a user with a unique key.
}

func (s *Services) Verify(resp http.ResponseWriter, req *http.Request) {
	// TODO: The verification endpoint that
}

type HealthCheckResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func (s *Services) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	// Assume healthy to start, then add in errors

	healthcheck := HealthCheckResponse{Message: "No issues", Errors: make([]string, 0)}
	err := s.db.Ping()
	if err != nil {
		healthcheck.Errors = append(healthcheck.Errors, err.Error())
	}
	if len(healthcheck.Errors) > 0 {
		healthcheck.Message = "Service is unhealthy, see `errors` for more information"
	}
	json.NewEncoder(resp).Encode(healthcheck)
}
