// Package models contains the shared API data that will be used by the server, templates, and cli in the code base.
package models

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
