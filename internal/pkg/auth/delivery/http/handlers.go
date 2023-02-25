package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/mailru/easyjson"
	"net/http"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase, or auth.OnlineRepo) *AuthHandler {
	return &AuthHandler{
		uc:     uc,
		online: or,
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.LoginUserIsValid(user) {
		middleware.Response(w, models.Forbidden, nil)
		return
	}
	token, status := h.uc.SignIn(user)
	if status != models.Okey {
		middleware.Response(w, models.Unauthed, nil)
	}
	if status == models.Okey {
		status = h.online.UserOn(user)
	}
	middleware.Response(w, status, models.TokenView{Token: token})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}
	accesToken, err := middleware.ExtractTokenMetadata(r, middleware.ExtractToken)
	if err != nil || accesToken == nil {
		middleware.Response(w, models.BadRequest, nil)
		return
	}

	user.Login = accesToken.Login

	status := h.online.UserOff(user)
	middleware.Response(w, status, nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.UserIsValid(user) {
		middleware.Response(w, models.BadRequest, nil)
		return
	}
	token, status := h.uc.SignUp(user)
	if status == models.Okey {
		userCopy := models.LoginUser{Login: user.Login, EncryptedPassword: user.EncryptedPassword}
		status = h.online.UserOn(userCopy)
	}

	if token == "" || status != models.Okey {
		middleware.Response(w, status, nil)
		return
	}
	middleware.Response(w, status, models.TokenView{Token: token})
}

func (h *AuthHandler) AuthStatus(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}

	user.Login = r.URL.Query().Get("user")
	jwtData, err := middleware.ExtractTokenMetadata(r, middleware.ExtractToken)

	if user.Login == "" || jwtData == nil {
		middleware.Response(w, models.BadRequest, nil)
		return
	}

	if err != nil || jwtData.Login != user.Login {
		middleware.Response(w, models.Unauthed, nil)
		return
	}

	status := h.online.IsAuthed(user)
	if !status {
		middleware.Response(w, models.Unauthed, nil)
		return
	}
	middleware.Response(w, models.Okey, nil)
}
