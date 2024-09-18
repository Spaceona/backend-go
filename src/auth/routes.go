package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type authDeviceResponse struct {
	Token string
}

type Route[T any] struct {
	Authenticate func(r *http.Request) (T, error)
	WriteToken   []func(w http.ResponseWriter, r *http.Request, token T)
	OnError      func(w http.ResponseWriter)
}

func (d Route[string]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, tokenErr := d.Authenticate(r)
	if tokenErr != nil {
		d.OnError(w)
		return
	}
	for _, writeToken := range d.WriteToken {
		writeToken(w, r, token)
	}
	return
}

func AuthHttpError(w http.ResponseWriter) {
	http.Error(w, "failed to authorize", http.StatusBadRequest)
}

func WriteTokenToAuthHeader(w http.ResponseWriter, r *http.Request, token string) {
	w.Header().Add("Authorization", "Bearer "+token)
}

func WriteTokenToBody(w http.ResponseWriter, r *http.Request, token string) {
	resStruct := authDeviceResponse{Token: token}
	resJson, jsonDecodeErr := json.Marshal(resStruct)
	if jsonDecodeErr != nil {
		slog.Error(jsonDecodeErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, writeErr := w.Write(resJson)
	if writeErr != nil {
		slog.Error(writeErr.Error())
		return
	}
}
