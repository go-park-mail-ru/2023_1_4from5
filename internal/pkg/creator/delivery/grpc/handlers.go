package grpcCreator

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
	"github.com/google/uuid"
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

func (h GrpcCreatorHandler) GetFeed(ctx context.Context, in *generatedCommon.UUIDMessage) (*generatedCreator.PostsMessage, error) {
	userID, err := uuid.Parse(in.Value)
	if err != nil {
		return &generatedCreator.PostsMessage{Error: err.Error()}, nil
	}
	feed, err := h.uc.GetFeed(ctx, userID)
	if err != nil {
		return &generatedCreator.PostsMessage{Error: err.Error()}, nil
	}

	var postsProto generatedCreator.PostsMessage
	for i, post := range feed {
		postsProto.Posts = append(postsProto.Posts, &generatedCreator.Post{
			Id:           post.Id.String(),
			CreatorID:    post.Creator.String(),
			Creation:     post.Creation.String(),
			CreatorName:  post.CreatorName,
			LikesCount:   post.LikesCount,
			CreatorPhoto: post.CreatorPhoto.String(),
			Title:        post.Title,
			Text:         post.Text,
			IsAvailable:  true,
			IsLiked:      post.IsLiked,
		})

		for _, attach := range post.Attachments {
			var attachsProto generatedCreator.Attachment
			attachsProto.ID = attach.Id.String()
			attachsProto.Type = attach.Type
			postsProto.Posts[i].PostAttachments = append(postsProto.Posts[i].PostAttachments, &attachsProto)
		}
		postsProto.Posts[i].Subscriptions = nil

	}
	postsProto.Error = ""

	return &postsProto, nil
}

func (h GrpcCreatorHandler) GetPage(ctx context.Context, in *generatedCreator.UserCreatorMessage) (*generatedCreator.CreatorPage, error) {
	return nil, nil
}

func (h GrpcCreatorHandler) UpdateCreatorData(ctx context.Context, in *generatedCreator.UpdateCreatorInfo) (*generatedCommon.Empty, error) {
	creatorID, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.UpdateCreatorData(ctx, models.UpdateCreatorInfo{
		Description: in.Description,
		CreatorName: in.CreatorName,
		CreatorID:   creatorID})
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}
