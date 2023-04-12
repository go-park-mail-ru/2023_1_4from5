package utils

import (
	"net/http"
)

func ResponseWithCSRF(w http.ResponseWriter, token string) {
	w.Header().Set("X-CSRF-Token", token)
	w.WriteHeader(http.StatusOK)
}
