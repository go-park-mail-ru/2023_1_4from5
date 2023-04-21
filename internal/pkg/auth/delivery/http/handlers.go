package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type AuthHandler struct {
	client generatedAuth.AuthServiceClient
	logger *zap.SugaredLogger
}

func NewAuthHandler(cl generatedAuth.AuthServiceClient, logger *zap.SugaredLogger) *AuthHandler {
	return &AuthHandler{
		client: cl,
		logger: logger,
	}
}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	user := models.LoginUser{}
	err := easyjson.UnmarshalFromReader(r.Body, &user)
	if err != nil || !(models.User{Login: user.Login, PasswordHash: user.PasswordHash}).UserAuthIsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	token, err := h.client.SignIn(r.Context(), &generatedAuth.LoginUser{
		Login:        user.Login,
		PasswordHash: user.PasswordHash,
	})
	if len(token.Error) != 0 {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Cookie(w, token.Cookie, "SSID")
	utils.Response(w, http.StatusOK, nil)
}

//func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
//	userData, err := token.ExtractJWTTokenMetadata(r)
//	if err != nil {
//		utils.Response(w, http.StatusBadRequest, nil)
//		return
//	}
//	if _, err := h.usecase.IncUserVersion(r.Context(), *userData); err != nil {
//		utils.Response(w, http.StatusInternalServerError, nil)
//		return
//	}
//	utils.Cookie(w, "", "SSID")
//	utils.Response(w, http.StatusOK, nil)
//}
//
//func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {
//	var user models.User
//	err := easyjson.UnmarshalFromReader(r.Body, &user)
//	if err != nil || !user.UserIsValid() {
//		utils.Response(w, http.StatusBadRequest, nil)
//		return
//	}
//
//	token, err := h.usecase.SignUp(r.Context(), user)
//	if token == "" {
//		if err == models.WrongData {
//			utils.Response(w, http.StatusConflict, nil)
//			return
//		}
//		utils.Response(w, http.StatusInternalServerError, nil)
//		return
//	}
//
//	utils.Cookie(w, token, "SSID")
//	utils.Response(w, http.StatusOK, nil)
//}
