package utils

import (
	"net/http"
	"time"
)

func Cookie(w http.ResponseWriter, token string) {
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   "sub-me.ru",
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(time.Hour * 24),
	}
	http.SetCookie(w, SSCookie)
}
