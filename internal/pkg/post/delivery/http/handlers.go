package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
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

// Создать массив uuid для attachs
// Создать attachs в файлах:
//Проверить наличие места подо все файлы сразу
//В цикле сохранить все файлы
// если ошибка - у нас есть индекс, до него проходимся и удаляем файлы
// Выполнить транзакцию

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	//TODO: проверка на соответствие userId и creatorId в кукеееееееее
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}
	/////////////////////////////////////////////////////////////////////////////
	r.Body = http.MaxBytesReader(w, r.Body, int64(models.MaxFormSize))
	err = r.ParseMultipartForm(models.MaxFormSize) // maxMemory
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if len(r.MultipartForm.File["attachments"]) > models.MaxFiles {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	for _, headers := range r.MultipartForm.File {
		for _, header := range headers {
			if header.Size > int64(models.MaxFileSize) {
				utils.Response(w, http.StatusBadRequest, nil)
				return
			}
		}
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
	///////////////////////////////////////////////////////////////////////////////

	postData.Attachments = make([]models.AttachmentData, len(postFilesTmp))
	for i, file := range postFilesTmp {
		tmpFile, err := file.Open()
		if err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
		buf, _ := io.ReadAll(tmpFile)

		if err = tmpFile.Close(); err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}

		if postData.Attachments[i].Data, err = file.Open(); err != nil {
			utils.Response(w, http.StatusBadRequest, nil)
			return
		}
		postData.Attachments[i].Type = http.DetectContentType(buf)
		postData.Attachments[i].Id = uuid.New()
	}

	if err = h.attachmentUsecase.CreateAttaches(postData.Attachments...); err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	} else if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	postData.Id = uuid.New()
	if err := h.usecase.CreatePost(r.Context(), postData); err != nil {
		_ = h.attachmentUsecase.DeleteAttaches(postData.Attachments...)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	utils.Response(w, http.StatusOK, postData)
}

func (h *PostHandler) AddLike(w http.ResponseWriter, r *http.Request) {
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.AddLike(r.Context(), userData.Id, like.PostID)
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
	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	var like models.Like
	err = easyjson.UnmarshalFromReader(r.Body, &like)
	if err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	like, err = h.usecase.RemoveLike(r.Context(), userData.Id, like.PostID)
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

	userData, err := token.ExtractJWTTokenMetadata(r)
	if err != nil {
		utils.Response(w, http.StatusUnauthorized, nil)
		return
	}

	if _, err := h.authUsecase.CheckUserVersion(r.Context(), *userData); err != nil {
		utils.Cookie(w, "", "SSID")
		utils.Response(w, http.StatusForbidden, nil)
		return
	}

	ok, err = h.usecase.IsPostOwner(r.Context(), (*userData).Id, postID)
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

	//err = h.attachmentUsecase.DeleteAttachsByPostID(postID) //DON't WORK
	if err == models.WrongData {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if err = h.usecase.DeletePost(r.Context(), postID); err != nil {
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	utils.Response(w, http.StatusOK, nil)
}
