package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
	"net/http"
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
	utils.Cookie(w, token)
	utils.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	_, err := jwt.ExtractTokenMetadata(r, jwt.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	utils.Cookie(w, "")
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
	if token == "" {
		if err == models.WrongData {
			utils.Response(w, http.StatusConflict, nil)
			return
		}
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Cookie(w, token)
	utils.Response(w, http.StatusOK, nil)
}
