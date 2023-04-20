package grpcAuth

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
)

//go:generate mockgen -source=auth_grpc.pb.go -destination=auth_grpc.go -package=grpc

type GrpcAuthHandler struct {
	uc auth.AuthUsecase
	generated.AuthServiceServer
}

func NewGrpcAuthHandler(uc auth.AuthUsecase) *GrpcAuthHandler {
	return &GrpcAuthHandler{
		uc: uc,
	}
}

//TODO:функционал
