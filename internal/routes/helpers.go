package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type APIError struct {
	StatusCode int      `json:"status_code"`
	Message    string   `json:"message"`
	Details    []string `json:"details"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error, details ...string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    err.Error(),
		Details:    details,
	}
}

func jsonError(err error) APIError {
	return APIError{
		StatusCode: http.StatusBadRequest,
		Message:    "Json could not be processed",
		Details:    []string{err.Error()},
	}
}

func forbiddenError() APIError {
	return APIError{
		StatusCode: http.StatusForbidden,
	}
}

func authenticationError(message string) APIError {
	return APIError{
		StatusCode: http.StatusUnauthorized,
		Message:    message,
	}
}

func internalServerError() APIError {
	return APIError{
		StatusCode: 500,
		Message:    "Internal Server Error",
		Details:    []string{"I am very sorry, an internal error has occurred with your request."},
	}
}

// Simple wrapper to reduce boilerplate over writing json in the api endpoints
func writeJSON(w http.ResponseWriter, status int, content any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(content)
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func Handler(handler APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {
				writeJSON(w, apiErr.StatusCode, apiErr)
			} else {
				// Catch all generic error
				content := map[string]any{
					"status_code": http.StatusInternalServerError,
					"message":     "Internal server error",
					"details":     "Please contact support regarding this issue.",
				}
				writeJSON(w, 500, content)
			}
			log.Printf("HTTP API error %v %v %v %v", "err", err.Error(), "path", r.URL.Path)
		}
	}
}
