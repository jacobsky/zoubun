package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	sqlc "zoubun/internal/db"

	"github.com/google/uuid"
)

type ErrorResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details"`
}

type Counter struct {
	Count int64 `json:"count"`
}

type Motd struct {
	Message string `json:"message"`
}

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

// General handler for errors
func processError(resp http.ResponseWriter, code int, message string, details ...string) {
	resp.WriteHeader(code)
	errResp := ErrorResponse{
		Message: message,
		Details: details,
	}
	json.NewEncoder(resp).Encode(errResp)
}

func dbError(resp http.ResponseWriter) {
	processError(resp, 500, "Database error has occurred", "Please contact the admin with the request ID for more information")
}

func (s *Services) MessageOfTheDay(resp http.ResponseWriter, req *http.Request) {
	message, err := s.queries.Motd(req.Context())
	if err != nil {
		log.Printf("Database Error %v", err.Error())
		resp.WriteHeader(503)
		errResp := ErrorResponse{
			Message: "An internal database error has occurred. If this issue persists, please notify the administrators",
		}
		json.NewEncoder(resp).Encode(errResp)
		return
	}

	motd := Motd{Message: message}
	json.NewEncoder(resp).Encode(motd)
}

func (s *Services) Count(resp http.ResponseWriter, req *http.Request) {
	userid := req.Header.Get("userid")
	if userid == "" {
		log.Print("/api/count somehow got a blank username from an authorization header")
		processError(resp, 500, "Internal Server Error", "")
		return
	}

	id, err := strconv.Atoi(userid)
	if err != nil {
		log.Print(err)
		resp.WriteHeader(500)
		return
	}

	count, err := s.queries.GetUserCounter(req.Context(), int32(id))
	if err != nil {
		log.Print(err)
		resp.WriteHeader(500)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(Counter{Count: count})
}

func (s *Services) Increment(resp http.ResponseWriter, req *http.Request) {
	userid := req.Header.Get("userid")
	if userid == "" {
		log.Print("/api/count somehow got a blank username from an authorization header")
		processError(resp, 500, "Internal Server Error", "")
		return
	}
	id, err := strconv.Atoi(userid)
	if err != nil {
		log.Print(err)
		resp.WriteHeader(500)
		return
	}

	count, err := s.queries.IncrementCounter(req.Context(), int32(id))
	if err != nil {
		log.Print(err)
		resp.WriteHeader(500)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json.NewEncoder(resp).Encode(Counter{Count: count})
}

type RegisterRequestBody struct {
	Username string `json:"username"`
}

type RegisterResponseBody struct {
	Username string `json:"username"`
	APIKey1  string `json:"apikey1"`
	APIKey2  string `json:"apikey2"`
	Message  string `json:"message"`
}

func (s *Services) Register(resp http.ResponseWriter, req *http.Request) {
	var body RegisterRequestBody
	json.NewDecoder(req.Body).Decode(&body)
	if len(body.Username) < 3 {
		log.Print("Empty username")
		processError(resp, 400, "Empty username", "Please select a username")
		return
	}
	exists, err := s.queries.UsernameExists(req.Context(), body.Username)
	if err != nil {
		fmt.Printf("Database Error has Occurred: %v", err)
		dbError(resp)
		return
	}

	if exists {
		log.Printf("Attempted to sign up with %v which already exists", body.Username)
		processError(resp, 400, "Username already exists", "Please choose another username and attempt to register again")
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Database Error has Occurred: %v", err)
		dbError(resp)
		return
	}

	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	userid, err := qtx.CreateUser(req.Context(), body.Username)
	if err != nil || !userid.Valid {
		log.Printf("userid [%v] is invalid and/or a database error occurred: %v", userid, err)
		dbError(resp)
	}
	// Create the API keys using the UUID generator (could be other methods, this is just simple)
	apikey1 := uuid.New().String()
	apikey2 := uuid.New().String()

	err = qtx.AddUserKey(req.Context(), sqlc.AddUserKeyParams{
		Userid:  int32(userid.Int32),
		Apikey1: apikey1,
		Apikey2: apikey2,
	})
	if err != nil {
		log.Printf("Database Error: %v", err)
		dbError(resp)
		return
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Database Error: %v", err)
		dbError(resp)
	}
	respBody := RegisterResponseBody{
		Username: body.Username,
		APIKey1:  apikey1,
		APIKey2:  apikey2,
		Message:  "Your key has been created, please be sure to save it.",
	}
	resp.WriteHeader(201)
	json.NewEncoder(resp).Encode(respBody)
}

type RotateKeyRequestBody struct {
	WhichKey int `json:"which_key"`
}

type RotateKeyResponseBody struct {
	WhichKey int    `json:"which_key"`
	NewKey   string `json:"new_key"`
}

func (s *Services) RotateKey(resp http.ResponseWriter, req *http.Request) {
	var body RotateKeyRequestBody
	// Extract the username from the header that should be added by middleware
	userid, err := strconv.Atoi(req.Header.Get("userId"))
	if err != nil {
		log.Printf("Fetching userid from headers resulted in %v", err)
		processError(resp, 500, "Internal server error", "Please contact the admin with the request ID")
		return
	}
	json.NewDecoder(req.Body).Decode(&body)
	if body.WhichKey != 1 && body.WhichKey != 2 {
		processError(resp,
			400,
			"Key does not exist",
			"Valid values for this request are 1 or 2 as there are only two keys per account",
		)
		return
	}
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Database Error has Occurred: %v", err.Error())
		dbError(resp)
		return
	}
	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	// Create the API keys using the UUID generator (could be other methods, this is just simple)
	newkey := uuid.New().String()

	switch body.WhichKey {
	case 1:
		_, err = qtx.RotateUserKey1(req.Context(), sqlc.RotateUserKey1Params{NewKey: newkey, Userid: int32(userid)})
	case 2:
		_, err = qtx.RotateUserKey2(req.Context(), sqlc.RotateUserKey2Params{NewKey: newkey, Userid: int32(userid)})
	}

	if err != nil {
		dbError(resp)
		return
	}
	respBody := RotateKeyResponseBody{
		WhichKey: body.WhichKey,
		NewKey:   newkey,
	}
	resp.WriteHeader(201)
	json.NewEncoder(resp).Encode(respBody)
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
