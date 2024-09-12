package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type authDeviceResponse struct {
	Token string
}

type Route[T any] struct {
	CanAuthenticate func(T) bool
}

func (d Route[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var body T
	jsonErr := json.NewDecoder(r.Body).Decode(&body)
	if jsonErr != nil {
		slog.Error(jsonErr.Error())
		http.Error(w, "could not authenticate device", http.StatusBadRequest)
		return
	}

	if d.CanAuthenticate(body) == false {
		http.Error(w, "could not authenticate device", http.StatusBadRequest)
		return
	}

	token, tokenErr := GenToken(body, 24*time.Hour)
	if tokenErr != nil {
		slog.Error(tokenErr.Error())
		http.Error(w, "could not authenticate device", http.StatusBadRequest)
		return
	}
	response := authDeviceResponse{token}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resJson, jsonEncodeErr := json.Marshal(response)
	if jsonEncodeErr != nil {
		slog.Error(jsonEncodeErr.Error())
		http.Error(w, "could not authenticate device", http.StatusBadRequest)
		return
	}
	_, writeErr := w.Write(resJson)
	if writeErr != nil {
		slog.Error(writeErr.Error())
		http.Error(w, "could not authenticate device", http.StatusBadRequest)
		return
	}
	return
}
