package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	"net/http"
)

type UserHandler struct {
	usecase user.UserUsecase
}

func NewUserHandler(uc user.UserUsecase) *UserHandler {
	return &UserHandler{
		usecase: uc,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userInfo, err := middleware.ExtractTokenMetadata(r, middleware.ExtractTokenFromCookie)
	if err != nil {
		middleware.Response(w, http.StatusUnauthorized, nil)
	}

	userProfile, err := h.usecase.GetProfile(*userInfo)
	if err != nil {
		middleware.Response(w, http.StatusInternalServerError, nil)
	}

	middleware.Response(w, http.StatusOK, userProfile)
}
