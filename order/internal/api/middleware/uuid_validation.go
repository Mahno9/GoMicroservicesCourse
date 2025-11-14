package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// UUIDValidationMiddleware validates UUID parameters in requests
func UUIDValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract UUID from path parameters
		pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

		// Check if this is an order endpoint with UUID
		if len(pathParts) >= 3 && pathParts[0] == "api" && pathParts[1] == "v1" && pathParts[2] == "order" {
			if len(pathParts) >= 4 {
				orderUUID := pathParts[3]
				if err := validateUUID(orderUUID); err != nil {
					http.Error(w, "Invalid order UUID format", http.StatusBadRequest)
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// validateUUID checks if the provided string is a valid UUID
func validateUUID(uuidStr string) error {
	_, err := uuid.Parse(uuidStr)
	return err
}
