package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/mailru/easyjson"
	"net/http"
)

type PostHandler struct {
	usecase           post.PostUsecase
	authUsecase       auth.AuthUsecase
	attachmentUsecase attachment.AttachmentUsecase
}

func NewPostHandler(uc post.PostUsecase, auc auth.AuthUsecase, attuc attachment.AttachmentUsecase) *PostHandler {
	return &PostHandler{
		usecase:           uc,
		authUsecase:       auc,
		attachmentUsecase: attuc,
	}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
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

	var postData models.PostCreationData
	if err = easyjson.UnmarshalFromReader(r.Body, &postData); err != nil || !postData.IsValid() {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postUUID := uuid.New()

	if err := h.attachmentUsecase.CreateAttach(postUUID, postData.Attachments...); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	postData.Id = postUUID
	if err := h.usecase.CreatePost(postData); err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, postUUID)
}

func (h *PostHandler) AddLike(w http.ResponseWriter, r *http.Request) {
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

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.AddLike(userData.Id, like.PostID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, like)
}

func (h *PostHandler) RemoveLike(w http.ResponseWriter, r *http.Request) {
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

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.RemoveLike(userData.Id, like.PostID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err == models.InternalError {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	utils.Response(w, http.StatusOK, like)
}
