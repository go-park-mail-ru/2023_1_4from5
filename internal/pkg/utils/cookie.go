package utils

import (
	"net/http"
	"os"
	"time"
)

func Cookie(w http.ResponseWriter, token string) {
	domain, _ := os.LookupEnv("DOMAIN") //TODO: обработать ошибку, ещё есть места где errors.NEw() прокидываются
	SSCookie := &http.Cookie{
		Name:   "SSID",
		Value:  token,
		Path:   "/",
		Domain: domain,
		// SameSite: 2,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(time.Hour * 24),
	}
	http.SetCookie(w, SSCookie)
}
