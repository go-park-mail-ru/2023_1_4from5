package utils

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

func Cookie(w http.ResponseWriter, token string) {
	domain, _ := os.LookupEnv("DOMAIN") //TODO: обработать ошибку, ещё есть места где errors.NEw() прокидываются
	fmt.Println(domain)
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
