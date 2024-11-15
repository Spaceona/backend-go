package auth

import (
	"fmt"
	"log/slog"
	"net/http"
	"spacesona-go-backend/logging"
	"strings"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("huh")
		//TODO validate token
		authHeader := r.Header.Get("Authorization")
		//todo get cookie
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
			slog.Error(validationErr.Error())
			logging.RequestsAuthenticatedFailed.Inc()
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}

		logging.RequestsAuthenticated.Inc()
		next.ServeHTTP(w, r)
	}
}
