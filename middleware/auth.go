package middleware

import (
	"context"
	"net/http"

	"github.com/nishanth-thoughtclan/student-api/utils"
)

// JWTAuthMiddleware verifies the JWT token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		token, err := utils.ValidateJWTToken(authHeader)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Attach claims to context
		ctx := context.WithValue(r.Context(), "claims", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
