package grpcUser

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	"github.com/google/uuid"
)

//go:generate mockgen -source=user_grpc.pb.go -destination=user_grpc.go -package=grpc
//go:generate protoc  --go_out=./generated --go-grpc_out=./generated --proto_path=. user.proto

type GrpcUserHandler struct {
	uc user.UserUsecase

	generatedUser.UserServiceServer
}

func NewGrpcUserHandler(uc user.UserUsecase) *GrpcUserHandler {
	return &GrpcUserHandler{
		uc: uc,
	}
}

func (h GrpcUserHandler) Follow(ctx context.Context, in *generatedUser.FollowMessage) (*generatedUser.Empty, error) {
	userId, err := uuid.Parse(in.UserID)
	creatorId, err := uuid.Parse(in.CreatorID)
	if err != nil {
		return &generatedUser.Empty{Error: err.Error()}, nil
	}
	err = h.uc.Follow(ctx, userId, creatorId)
	if err != nil {
		return &generatedUser.Empty{Error: err.Error()}, nil
	}
	return &generatedUser.Empty{Error: ""}, nil
}
