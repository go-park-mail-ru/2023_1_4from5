package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/jwt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"io"
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
	//TODO: проверка на соответствие userId и creatorId
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

	err = r.ParseMultipartForm(32 << 20) // maxMemory
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	postValues := r.MultipartForm.Value
	postFilesTmp := r.MultipartForm.File["attachments"]

	var postData models.PostCreationData

	if postData.Creator, err = uuid.Parse(postValues["creator"][0]); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	postData.Text = postValues["text"][0]

	postData.Title = postValues["title"][0]

	tmpSubs := postValues["subscriptions"]
	postData.AvailableSubscriptions = make([]uuid.UUID, len(tmpSubs))
	for i, sub := range tmpSubs {
		if postData.AvailableSubscriptions[i], err = uuid.Parse(sub); err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
	}

	postData.Id, err = h.usecase.CreatePost(postData)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	//TODO: в случае ошибки удалить пост

	attachments := make([]models.AttachmentData, len(postFilesTmp))

	for i, file := range postFilesTmp {
		attachments[i].Data, err = file.Open()
		if err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}

		buf, _ := io.ReadAll(attachments[i].Data)
		attachments[i].Data.Close()
		if attachments[i].Data, err = file.Open(); err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
		attachments[i].Type = http.DetectContentType(buf)
	}

	postData.Attachments = make([]uuid.UUID, len(attachments))
	if postData.Attachments, err = h.attachmentUsecase.CreateAttachs(postData.Id, attachments...); err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	} else if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, postData)
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

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDtmp, ok := mux.Vars(r)["post-uuid"]

	if !ok {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	postID, err := uuid.Parse(postIDtmp)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

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

	ok, err = h.usecase.IsPostOwner((*userData).Id, postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if !ok {
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	err = h.attachmentUsecase.DeleteAttachsByPostID(postID)
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err = h.usecase.DeletePost(postID); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}
