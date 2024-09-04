package info

import (
	"net/http"
)

func APIInfoRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, writeErr := w.Write([]byte("Spaceona backend v2.0"))
	if writeErr != nil {
		http.Error(w, writeErr.Error(), http.StatusInternalServerError)
		return
	}
}
