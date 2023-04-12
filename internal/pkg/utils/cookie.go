package utils

import (
	"net/http"
	"os"
	"time"
)

func Cookie(w http.ResponseWriter, token, name string) {
	domain, _ := os.LookupEnv("DOMAIN")
	SSCookie := &http.Cookie{
		Name:     name,
		Value:    token,
		Path:     "/",
		Domain:   domain,
		SameSite: http.SameSiteStrictMode,
		HttpOnly: true,
		Expires:  time.Now().UTC().Add(time.Hour * 24),
	}
	http.SetCookie(w, SSCookie)
}
