package routes

import (
	"encoding/json"
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
			resp.WriteHeader(401)
			json.NewEncoder(resp).Encode(AuthorizationError{
				Kind:    "Unauthorized",
				Details: "No authorization header was found. Please ensure that the API key is included in the `zoubun-api-key` header",
			})
			return
		}
		// TODO: Plumb it into a DB query
		id, err := s.queries.GetUserIdFromAuth(req.Context(), key)
		if err != nil {
			log.Printf("Database Error: %v", err)
			resp.WriteHeader(500)
			json.NewEncoder(resp).Encode(AuthorizationError{
				Kind:    "Interal Authorization Error",
				Details: "An error occurred when attempting to authorize your request.",
			})
			return
		}

		if !id.Valid {
			resp.WriteHeader(403)
			json.NewEncoder(resp).Encode(AuthorizationError{
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
