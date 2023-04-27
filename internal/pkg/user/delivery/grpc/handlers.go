package grpcUser

import (
	"context"
	generatedCommon "github.com/go-park-mail-ru/2023_1_4from5/internal/models/proto"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/google/uuid"
)

//go:generate mockgen -source=user_grpc.pb.go -destination=user_grpc.go -package=grpc

type GrpcUserHandler struct {
	uc user.UserUsecase

	generatedUser.UserServiceServer
}

func NewGrpcUserHandler(uc user.UserUsecase) *GrpcUserHandler {
	return &GrpcUserHandler{
		uc: uc,
	}
}

func (h GrpcUserHandler) Follow(ctx context.Context, in *generatedUser.FollowMessage) (*generatedCommon.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	err = h.uc.Follow(ctx, userId, creatorId)
	if err != nil {
		return &generatedCommon.Empty{Error: err.Error()}, nil
	}
	return &generatedCommon.Empty{Error: ""}, nil
}
