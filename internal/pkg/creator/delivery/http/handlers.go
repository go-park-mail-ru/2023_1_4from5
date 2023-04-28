package http

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/token"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"net/http"
)

type CreatorHandler struct {
	creatorClient generatedCreator.CreatorServiceClient
	usecase       creator.CreatorUsecase
	authClient    generatedAuth.AuthServiceClient
	postUsecase   post.PostUsecase
	logger        *zap.SugaredLogger
}

func NewCreatorHandler(uc creator.CreatorUsecase, puc post.PostUsecase, creatorClient generatedCreator.CreatorServiceClient, authClient generatedAuth.AuthServiceClient, logger *zap.SugaredLogger) *CreatorHandler {
	return &CreatorHandler{
		usecase:       uc,
		creatorClient: creatorClient,
		authClient:    authClient,
		postUsecase:   puc,
		logger:        logger,
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

	out, err := h.creatorClient.FindCreators(r.Context(), &generatedCreator.KeywordMessage{Keyword: keyword})
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	if out.Error != "" {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}

	creators := make([]models.Creator, len(out.Creators))

	for i, v := range out.Creators {
		creatorId, err := uuid.Parse(v.Id)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		userId, err := uuid.Parse(v.UserID)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creatorPhoto, err := uuid.Parse(v.CreatorPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		coverPhoto, err := uuid.Parse(v.CoverPhoto)
		if err != nil {
			utils.Response(w, http.StatusInternalServerError, nil)
			return
		}
		creators[i] = models.Creator{
			Id:             creatorId,
			UserId:         userId,
			Name:           v.CreatorName,
			CoverPhoto:     coverPhoto,
			ProfilePhoto:   creatorPhoto,
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		}

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

	creatorId, err := uuid.Parse(creatorUUID)
	if err != nil {
		h.logger.Error(err)
		utils.Response(w, http.StatusBadRequest, nil)
		return
	}

	creatorPage, err := h.usecase.GetPage(r.Context(), tmpUserInfo.Id, creatorId)
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

	uv, err := h.authClient.CheckUserVersion(r.Context(), &generatedAuth.AccessDetails{
		Login:       userDataJWT.Login,
		Id:          userDataJWT.Id.String(),
		UserVersion: int64(userDataJWT.UserVersion),
	})
	if err != nil {
		utils.Response(w, http.StatusInternalServerError, nil)
		return
	}
	if len(uv.Error) != 0 {
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
