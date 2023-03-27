package http

import (
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
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
	userInfo, err := token.ExtractJWTTokenMetadata(r)
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
	userInfo, err := token.ExtractJWTTokenMetadata(r)
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
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)

	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(*userDataJWT); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if r.Method == http.MethodGet {
		tokenCSRF, err := token.GetCSRFToken(models.User{Login: userDataJWT.Login, Id: userDataJWT.Id, UserVersion: userDataJWT.UserVersion})
		if err != nil {
			utils.Response(w, http.StatusUnauthorized, nil)
			return
		}
		utils.Cookie(w, tokenCSRF, "X-CSRF-Token")
		return
	}
	// check CSRF token
	userDataCSRF, err := token.ExtractCSRFTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	if *userDataCSRF != *userDataJWT {
		fmt.Println(userDataJWT, " ", userDataCSRF)
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = r.ParseMultipartForm(4 << 20) // maxMemory
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	file, fileTmp, err := r.FormFile("upload")
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	// проверка типа файла
	buf, _ := io.ReadAll(file)
	file.Close()
	if file, err = fileTmp.Open(); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if http.DetectContentType(buf) != "image/jpeg" && http.DetectContentType(buf) != "image/png" {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	defer file.Close()

	var oldName uuid.UUID
	oldName, err = uuid.Parse(r.PostFormValue("path"))
	if err != nil {
		fmt.Println(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if oldName != uuid.Nil {
		err = os.Remove(fmt.Sprintf("%s%s.jpg", models.FolderPath, oldName.String()))
		if err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
		}
	}

	name, err := h.usecase.UpdatePhoto(*userDataJWT)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	f, err := os.Create(fmt.Sprintf("%s%s.jpg", models.FolderPath, name.String()))
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
