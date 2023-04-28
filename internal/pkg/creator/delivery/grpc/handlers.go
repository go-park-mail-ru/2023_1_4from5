package grpcCreator

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"golang.org/x/net/context"
)

type GrpcCreatorHandler struct {
	uc  creator.CreatorUsecase
	puc post.PostUsecase
	generatedCreator.CreatorServiceServer
}

func NewGrpcCreatorHandler(uc creator.CreatorUsecase, puc post.PostUsecase) *GrpcCreatorHandler {
	return &GrpcCreatorHandler{
		uc:  uc,
		puc: puc,
	}
}

func (h GrpcCreatorHandler) FindCreators(ctx context.Context, in *generatedCreator.KeywordMessage) (*generatedCreator.CreatorsMessage, error) {
	creators, err := h.uc.FindCreators(ctx, in.Keyword)
	if err != nil {
		return &generatedCreator.CreatorsMessage{Error: err.Error()}, nil
	}
	var creatorsMessage generatedCreator.CreatorsMessage
	for _, v := range creators {
		creatorsMessage.Creators = append(creatorsMessage.Creators, &generatedCreator.Creator{
			Id:             v.Id.String(),
			UserID:         v.UserId.String(),
			CreatorName:    v.Name,
			CreatorPhoto:   v.ProfilePhoto.String(),
			CoverPhoto:     v.CoverPhoto.String(),
			FollowersCount: v.FollowersCount,
			Description:    v.Description,
			PostsCount:     v.PostsCount,
		})
	}
	creatorsMessage.Error = ""

	return &creatorsMessage, nil
}

func (h GrpcCreatorHandler) GetPage(ctx context.Context, in *generatedCreator.UserCreatorMessage) (*generatedCreator.CreatorPage, error) {
	return nil, nil
}
