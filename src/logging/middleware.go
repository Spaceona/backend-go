package logging

import (
	"log/slog"
	"net/http"
)

func Middleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Got request", "Remote Address", r.RemoteAddr, "Route", r.RequestURI, "Method", r.Method)
		RequestsReceived.Inc()
		next.ServeHTTP(w, r)

	}
}
