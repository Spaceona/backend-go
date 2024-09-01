package auth

import (
	"net/http"
	"spacesona-go-backend/logging"
	"strings"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//TODO validate token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			logging.RequestsAuthenticatedFailed.Inc()
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		token := strings.Split(authHeader, " ")
		if len(token) < 2 {
			logging.RequestsAuthenticatedFailed.Inc()
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		validationErr := ValidateToken(token[1])
		if validationErr != nil {
			logging.RequestsAuthenticatedFailed.Inc()
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		
		logging.RequestsAuthenticated.Inc()
		next.ServeHTTP(w, r)
	}
}
