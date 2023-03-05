package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
	"net/http"
	"os"
	"time"
)

type AuthHandler struct {
	usecase auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: uc,
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}
	url, _ := os.LookupEnv("URL")
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.UserIsValid(models.User{Login: user.Login, PasswordHash: user.PasswordHash}) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token, status := h.usecase.SignIn(user)
	if status != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   url,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, SSCookie)
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.ExtractTokenMetadata(r, middleware.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	url, _ := os.LookupEnv("URL")
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    "",
		Path:     "/",
		Domain:   url,
		HttpOnly: true,
		Expires:  time.Now(),
	}
	http.SetCookie(w, SSCookie)
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.UserIsValid(user) {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token, status := h.usecase.SignUp(user)
	if token == "" || token == "no secret key" { //TODO в константу
		if status == models.ConflictData {
			utils.Response(w, http.StatusConflict, nil)
			return
		}
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	url, _ := os.LookupEnv("URL")
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   url,
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, SSCookie)
	w.WriteHeader(http.StatusOK)
}
