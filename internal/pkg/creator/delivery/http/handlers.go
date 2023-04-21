package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"net/http"
)

type CreatorHandler struct {
	usecase     creator.CreatorUsecase
	authUsecase auth.AuthUsecase
	postUsecase post.PostUsecase
}

func NewCreatorHandler(uc creator.CreatorUsecase, auc auth.AuthUsecase, puc post.PostUsecase) *CreatorHandler {
	return &CreatorHandler{
		usecase:     uc,
		authUsecase: auc,
		postUsecase: puc,
	}
}

func (h *CreatorHandler) GetAllCreators(w http.ResponseWriter, r *http.Request) {
	creators, err := h.usecase.GetAllCreators(r.Context())
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	for i := range creators {
		creators[i].Sanitize()
	}
	utils.Response(w, http.StatusOK, creators)
}

func (h *CreatorHandler) FindCreator(w http.ResponseWriter, r *http.Request) {
	keyword, ok := mux.Vars(r)["keyword"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creators, err := h.usecase.FindCreators(r.Context(), keyword)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	for i, _ := range creators {
		creators[i].Sanitize()
	}
	utils.Response(w, http.StatusOK, creators)
}

func (h *CreatorHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	creatorUUID, ok := mux.Vars(r)["creator-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	userInfo := models.AccessDetails{}
	tmpUserInfo, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		tmpUserInfo = &userInfo
	}
	creatorPage, err := h.usecase.GetPage(r.Context(), tmpUserInfo, creatorUUID)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	creatorPage.Sanitize()

	utils.Response(w, http.StatusOK, creatorPage)
}

func (h *CreatorHandler) CreateAim(w http.ResponseWriter, r *http.Request) {
	userDataJWT, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userDataJWT); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	aimInfo := models.Aim{}
	err = easyjson.UnmarshalFromReader(r.Body, &aimInfo)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if len(aimInfo.Description) > 100 {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	isCreator, err := h.postUsecase.IsCreator(r.Context(), userDataJWT.Id, aimInfo.Creator)

	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !isCreator {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = h.usecase.CreateAim(r.Context(), aimInfo)
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, nil)
}
