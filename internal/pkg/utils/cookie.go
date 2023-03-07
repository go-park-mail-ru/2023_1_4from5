package utils

import (
	"net/http"
	"time"
)

func Cookie(w http.ResponseWriter, url string, token string) {
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   url,
		HttpOnly: true,
		Expires:  time.Now(),
	}
	http.SetCookie(w, SSCookie)
}
