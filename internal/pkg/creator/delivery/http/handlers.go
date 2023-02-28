package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
)

type CreatorHandler struct {
	usecase creator.CreatorUsecase
}

func NewCreatorHandler(uc creator.CreatorUsecase) *CreatorHandler {
	return &CreatorHandler{
		usecase: uc,
	}
}

//func (h *CreatorHandler) GetPage(w http.ResponseWriter, r *http.Request) {
//	// TODO продумать отображение страницы автора в зависимости от уровня подписки

//}
