package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type CreatorHandler struct {
	usecase creator.CreatorUsecase
}

func NewCreatorHandler(uc creator.CreatorUsecase) *CreatorHandler {
	return &CreatorHandler{
		usecase: uc,
	}
}

func (h *CreatorHandler) GetPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	creatorUUID, ok := vars["creator-uuid"]
	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	userInfo, err := middleware.ExtractTokenMetadata(r, middleware.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}
	creatorPage, err := h.usecase.GetPage(*userInfo, creatorUUID)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, creatorPage)
}
