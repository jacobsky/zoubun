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

func (s *Services) MessageOfTheDay(resp http.ResponseWriter, req *http.Request) error {
	message, err := s.queries.Motd(req.Context())
	if err != nil {
		// log.Printf("Database Error %v", err.Error())
		return internalServerError()
	}

	motd := Motd{Message: message}
	return writeJSON(resp, http.StatusOK, motd)
}

func (s *Services) Count(resp http.ResponseWriter, req *http.Request) error {
	userid := req.Header.Get("userid")
	if userid == "" {
		log.Print("/api/count somehow got a blank username from an authorization header")
		return internalServerError()
	}

	id, err := strconv.Atoi(userid)
	if err != nil {
		return internalServerError()
	}

	count, err := s.queries.GetUserCounter(req.Context(), int32(id))
	if err != nil {
		return internalServerError()
	}

	return writeJSON(resp, http.StatusOK, Counter{Count: count})
}

func (s *Services) Increment(resp http.ResponseWriter, req *http.Request) error {
	userid := req.Header.Get("userid")
	if userid == "" {
		log.Print("/api/count somehow got a blank username from an authorization header")
		return internalServerError()
	}
	id, err := strconv.Atoi(userid)
	if err != nil {
		return internalServerError()
	}

	count, err := s.queries.IncrementCounter(req.Context(), int32(id))
	if err != nil {
		return internalServerError()
	}

	return writeJSON(resp, http.StatusOK, Counter{Count: count})
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

func (s *Services) Register(resp http.ResponseWriter, req *http.Request) error {
	var body RegisterRequestBody
	json.NewDecoder(req.Body).Decode(&body)
	if len(body.Username) < 3 {
		log.Print("Empty username")
		return NewAPIError(400, fmt.Errorf("empty username"), "Please select a username")
	}
	exists, err := s.queries.UsernameExists(req.Context(), body.Username)
	if err != nil {
		fmt.Printf("Database Error has Occurred: %v", err)
		return internalServerError()
	}

	if exists {
		log.Printf("Attempted to sign up with %v which already exists", body.Username)
		return NewAPIError(400, fmt.Errorf("Username already exists"), "Your desired username is already in use. Please select another name and register again")
	}
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Database Error has Occurred: %v", err)
		return internalServerError()
	}

	defer tx.Rollback()
	qtx := s.queries.WithTx(tx)
	userid, err := qtx.CreateUser(req.Context(), body.Username)
	if err != nil || !userid.Valid {
		log.Printf("userid [%v] is invalid and/or a database error occurred: %v", userid, err)
		return internalServerError()
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
		return internalServerError()
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Database Error: %v", err)
		return internalServerError()
	}
	respBody := RegisterResponseBody{
		Username: body.Username,
		APIKey1:  apikey1,
		APIKey2:  apikey2,
		Message:  "Your key has been created, please be sure to save it.",
	}

	return writeJSON(resp, http.StatusCreated, respBody)
}

type RotateKeyRequestBody struct {
	WhichKey int `json:"which_key"`
}

type RotateKeyResponseBody struct {
	WhichKey int    `json:"which_key"`
	NewKey   string `json:"new_key"`
}

func (s *Services) RotateKey(resp http.ResponseWriter, req *http.Request) error {
	var body RotateKeyRequestBody
	// Extract the username from the header that should be added by middleware
	userid, err := strconv.Atoi(req.Header.Get("userId"))
	if err != nil {
		log.Printf("Fetching userid from headers resulted in %v", err)
		return internalServerError()
	}
	json.NewDecoder(req.Body).Decode(&body)
	if body.WhichKey != 1 && body.WhichKey != 2 {
		return NewAPIError(
			400,
			fmt.Errorf("Key does not exist"),
			"Valid values for this request are 1 or 2 as there are only two keys per account",
		)
	}
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Database Error has Occurred: %v", err.Error())
		return internalServerError()
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
		return internalServerError()
	}
	respBody := RotateKeyResponseBody{
		WhichKey: body.WhichKey,
		NewKey:   newkey,
	}
	return writeJSON(resp, http.StatusCreated, respBody)
}

func (s *Services) Verify(resp http.ResponseWriter, req *http.Request) {
	// TODO: The verification endpoint that
}

type HealthCheckResponse struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func (s *Services) HealthCheck(resp http.ResponseWriter, req *http.Request) error {
	// Assume healthy to start, then add in errors

	healthcheck := HealthCheckResponse{Message: "No issues", Errors: make([]string, 0)}
	err := s.db.Ping()
	if err != nil {
		healthcheck.Errors = append(healthcheck.Errors, err.Error())
	}
	if len(healthcheck.Errors) > 0 {
		healthcheck.Message = "Service is unhealthy, see `errors` for more information"
	}
	return writeJSON(resp, http.StatusOK, healthcheck)
}
