package http

// TODO Add domain
import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/mailru/easyjson"
	"net/http"
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
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.LoginUserIsValid(user) {
		middleware.Response(w, http.StatusForbidden, nil)
		return
	}

	token, status := h.usecase.SignIn(user)

	if status != http.StatusOK {
		middleware.Response(w, http.StatusUnauthorized, nil)
		return
	}
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   "???????????????",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}
	http.SetCookie(w, SSCookie)
	middleware.Response(w, status, nil)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    "",
		Path:     "/",
		Domain:   "???????????????",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, SSCookie)
	middleware.Response(w, http.StatusOK, nil)
}

func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !middleware.UserIsValid(user) {
		middleware.Response(w, http.StatusBadRequest, nil)
		return
	}

	token, status := h.usecase.SignUp(user)

	if token == "" || token == "no secret key" { //TODO в константу
		middleware.Response(w, status, nil)
		return
	}

	SSCookie := &http.Cookie{
		Name:     "SSID",
		Value:    token,
		Path:     "/",
		Domain:   "?????????",
		HttpOnly: true,
		Expires:  time.Now().Add(time.Hour * 24),
	}

	http.SetCookie(w, SSCookie)
	w.WriteHeader(http.StatusOK)
}
