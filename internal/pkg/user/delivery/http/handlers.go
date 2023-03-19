package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"io"
	"net/http"
	"os"
)

type UserHandler struct {
	usecase     user.UserUsecase
	authUsecase auth.AuthUsecase
}

func NewUserHandler(uc user.UserUsecase, auc auth.AuthUsecase) *UserHandler {
	return &UserHandler{
		usecase:     uc,
		authUsecase: auc,
	}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userInfo, err := jwt.ExtractTokenMetadata(r, jwt.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	userProfile, err := h.usecase.GetProfile(*userInfo)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.NotFound {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, userProfile)
}

func (h *UserHandler) GetHomePage(w http.ResponseWriter, r *http.Request) {
	userInfo, err := jwt.ExtractTokenMetadata(r, jwt.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	homePage, err := h.usecase.GetHomePage(*userInfo)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.NotFound {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, homePage)
}

func (h *UserHandler) UpdateProfilePhoto(w http.ResponseWriter, r *http.Request) {
	userData, err := jwt.ExtractTokenMetadata(r, jwt.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(*userData); err != nil {
		utils.Cookie(w, "")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = r.ParseMultipartForm(4 << 20) // maxMemory
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	// удаляем старое фото с сервера
	// если фото не было, придёт uuid.Nil
	var oldPath uuid.UUID
	oldPath, err = uuid.Parse(r.PostFormValue("path"))
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if oldPath != uuid.Nil {
		err = os.Remove(fmt.Sprintf("/home/ubuntu/frontend/2023_1_4from5/public/%s.jpg", oldPath.String()))
		if err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
		}
	}

	file, _, err := r.FormFile("upload")
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	defer file.Close()
	path, err := h.usecase.UpdatePhoto(*userData)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	f, err := os.Create(fmt.Sprintf("/home/ubuntu/frontend/2023_1_4from5/public/%s.jpg", path.String()))
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	utils.Response(w, http.StatusOK, nil)
}
