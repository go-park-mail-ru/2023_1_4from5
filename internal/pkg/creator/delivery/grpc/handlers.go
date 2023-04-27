package grpcCreator

import (
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post"
)

type GrpcCreatorHandler struct {
	uc creator.CreatorUsecase
	pc post.PostUsecase
	generatedCreator.CreatorServiceServer
}

func NewGrpcCreatorHandler(uc creator.CreatorUsecase, pc post.PostUsecase) *GrpcCreatorHandler {
	return &GrpcCreatorHandler{
		uc: uc,
		pc: pc,
	}
}
