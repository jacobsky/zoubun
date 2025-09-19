// Package models contains the shared API data that will be used by the server, templates, and cli in the code base.
package api

import (
	"encoding/json"
	"net/http"
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

func Index(resp http.ResponseWriter, req *http.Request) {
	motd := Motd{Message: "皆さん、/incrementや/countのエンドポイントで増分してみよう～"}
	resp.Header().Set("ContentType", "text/html; charset=utf-8")
	json.NewEncoder(resp).Encode(motd)
}

func Count(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func Increment(resp http.ResponseWriter, req *http.Request) {
	currentCount.Count++
	resp.Header().Set("ContentType", "application/json")
	json.NewEncoder(resp).Encode(currentCount)
}

func Register(resp http.ResponseWriter, req *http.Request) {
	// TODO: This function will be used to register a user with a unique key.
}

func Verify(resp http.ResponseWriter, req *http.Request) {
	// TODO: The verification endpoint that
}
