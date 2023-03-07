package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
	"net/http"
	"os"
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
	url, flag := os.LookupEnv("URL")
	if !flag {
		//TODO
	}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !(models.User{Login: user.Login, PasswordHash: user.PasswordHash}).UserIsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token, err := h.usecase.SignIn(user)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	utils.Cookie(w, url, token)
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := middleware.ExtractTokenMetadata(r, middleware.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	url, _ := os.LookupEnv("URL")
	utils.Cookie(w, url, "")
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !user.UserIsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	token, err := h.usecase.SignUp(user)
	if token == "" || token == "no secret key" { //TODO в константу
		if err == models.ConflictData {
			utils.Response(w, http.StatusConflict, nil)
			return
		}
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	url, _ := os.LookupEnv("URL")

	utils.Cookie(w, url, token)
	utils.Response(w, http.StatusOK, nil)
}
