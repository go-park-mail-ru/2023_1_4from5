package main

import (
	"database/sql"
	grpcUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	userRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/repo"
	userUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}
func run() error {

	logger := utils.FileLogger("logUser.txt")

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			log.Print(err)
		}
	}(logger)

	zapSugar := logger.Sugar()

	str, err := utils.GetConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	userRepo := userRepository.NewUserRepo(db, zapSugar)
	userUse := userUsecase.NewUserUsecase(userRepo, zapSugar)
	service := grpcUser.NewGrpcUserHandler(userUse)

	srv, ok := net.Listen("tcp", ":8020") //TODO:разобраться с портами
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	generatedUser.RegisterUserServiceServer(server, service)

	log.Print("user running on: ", srv.Addr())
	return server.Serve(srv)
}
