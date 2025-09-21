package routes

import (
	"encoding/json"
	"net/http"

	db "zoubun/internal/db"
)

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
	db *db.Queries
}

func NewServices(db *db.Queries) *Services {
	return &Services{db}
}

func (s *Services) Index(resp http.ResponseWriter, req *http.Request) {
	motd := Motd{Message: "皆さん、/incrementや/countのエンドポイントで増分してみよう～"}
	resp.Header().Set("ContentType", "text/html; charset=utf-8")
	json.NewEncoder(resp).Encode(motd)
}

func (s *Services) Count(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func (s *Services) Increment(resp http.ResponseWriter, req *http.Request) {
	currentCount.Count++
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func (s *Services) Register(resp http.ResponseWriter, req *http.Request) {
	// TODO: This function will be used to register a user with a unique key.
}

func (s *Services) Verify(resp http.ResponseWriter, req *http.Request) {
	// TODO: The verification endpoint that
}

func (s *Services) HealthCheck(resp http.ResponseWriter, req *http.Request) {
	// TODO: Create health check endpoint that pings the main service.
}
