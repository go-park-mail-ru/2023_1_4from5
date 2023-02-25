package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
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
	var user models.LoginUser
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.LoginUserIsValid(user) {
		middleware.Response(w, http.StatusForbidden, nil)
		return
	}
	token, status := h.usecase.SignIn(user)
	if status != http.StatusOK {
		middleware.Response(w, http.StatusUnauthorized, nil)
	}

	middleware.Response(w, status, models.TokenView{Token: token})
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.UserIsValid(user) {
		middleware.Response(w, http.StatusBadRequest, nil)
		return
	}
	token, status := h.usecase.SignUp(user)

	if token == "" || status != http.StatusOK {
		middleware.Response(w, status, nil)
		return
	}
	middleware.Response(w, status, models.TokenView{Token: token})
}
