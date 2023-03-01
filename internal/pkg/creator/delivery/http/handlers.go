package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/mailru/easyjson"
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
	var creatorInfo models.Creator
	err := easyjson.UnmarshalFromReader(r.Body, &creatorInfo)
	if err != nil {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	userInfo, err := middleware.ExtractTokenMetadata(r, middleware.ExtractTokenFromCookie)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
	}

	creatorPage, err := h.usecase.GetPage(*userInfo, creatorInfo)
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
	}

	utils.Response(w, http.StatusOK, creatorPage)
}
