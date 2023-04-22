package main

import (
	"database/sql"
	grpcAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	authRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
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

	logger := utils.FileLogger("log.txt")

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

	tokenGenerator := authUsecase.NewTokenator()
	encryptor, err := authUsecase.NewEncryptor()
	if err != nil {
		return err
	}

	authRepo := authRepository.NewAuthRepo(db, zapSugar)
	authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator, encryptor, zapSugar)
	service := grpcAuth.NewGrpcAuthHandler(authUse)

	srv, ok := net.Listen("tcp", ":8010") //TODO:разобраться с портами
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	server := grpc.NewServer()

	generatedAuth.RegisterAuthServiceServer(server, service)

	log.Print("auth running on: ", srv.Addr())
	return server.Serve(srv)
}
