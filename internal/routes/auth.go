package routes

import (
	"log"
	"net/http"
	"strconv"
)

type AuthorizationError struct {
	Kind    string `json:"kind"`
	Details string `json:"details"`
}

func (s *Services) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// If the key is valid, then serve the request
		key := req.Header.Get("zoubun-api-key")
		if key == "" {
			// Unique case of ignoring because it is most _certainly_ a user error.
			err := authenticationError("No authorization header was found. Please ensure that the API key is included in the `zoubun-api-key` header")
			log.Printf("Auth Error: %v", err.Error())
			return
		}

		id, err := s.queries.GetUserIdFromAuth(req.Context(), key)
		if err != nil {
			log.Printf("Database Error: %v", err)
			writeJSON(resp, http.StatusInternalServerError, internalServerError())
			return
		}

		if !id.Valid {
			writeJSON(resp, http.StatusForbidden, AuthorizationError{
				Kind:    "Forbidden",
				Details: "The `zoubun-api-key` header is invalid or incorrect",
			})
			return
		}
		// Add the Userid to the request header so that it's easier to track
		req.Header.Add("userid", strconv.Itoa(int(id.Int32)))
		next.ServeHTTP(resp, req)
		// Otherwise return no
	})
}
