package main

import (
	"database/sql"
	"fmt"
	attachmentRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment/repo"
	attachmentUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment/usecase"
	commentRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/comment/repo"
	commentUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/comment/usecase"
	grpcCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	creatorRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/repo"
	creatorUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	postRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/repo"
	postUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/usecase"
	subscriptionRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription/repo"
	subscriptionUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
}

func run() error {
	logger, err := utils.FileLogger("/var/log/creator_app.log")
	if err != nil {
		return err
	}

	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Print(err)
		}
	}(logger)

	zapSugar := logger.Sugar()

	utils.Init()
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

	subscriptionRepo := subscriptionRepository.NewSubscriptionRepo(db, zapSugar)
	subscriptionUse := subscriptionUsecase.NewSubscriptionUsecase(subscriptionRepo, zapSugar)

	postRepo := postRepository.NewPostRepo(db, zapSugar)
	postUse := postUsecase.NewPostUsecase(postRepo, zapSugar)

	attachmentRepo := attachmentRepository.NewAttachmentRepo(db, zapSugar)
	attachmentUse := attachmentUsecase.NewAttachmentUsecase(attachmentRepo, zapSugar)

	creatorRepo := creatorRepository.NewCreatorRepo(db, zapSugar)
	creatorUse := creatorUsecase.NewCreatorUsecase(creatorRepo, zapSugar)

	commentRepo := commentRepository.NewCommentRepo(db, zapSugar)
	commentUse := commentUsecase.NewCommentUsecase(commentRepo, zapSugar)

	service := grpcCreator.NewGrpcCreatorHandler(creatorUse, postUse, attachmentUse, subscriptionUse, commentUse)

	srv, ok := net.Listen("tcp", ":8030")
	if ok != nil {
		log.Fatalln("can't listen port", err)
	}

	metricsMw := middleware.NewMetricsMiddleware()
	metricsMw.Register(middleware.ServiceCreatorName)

	server := grpc.NewServer(grpc.UnaryInterceptor(metricsMw.ServerMetricsInterceptor))

	generatedCreator.RegisterCreatorServiceServer(server, service)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: ":8031"}

	go func() {
		err := httpSrv.ListenAndServe()
		if err != nil {
			fmt.Print(err)
		}
	}()

	fmt.Print("creator running on: ", srv.Addr())
	return server.Serve(srv)
}
