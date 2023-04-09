package main

import (
	"database/sql"
	attachmentRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment/repo"
	attachmentUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/attachment/usecase"
	authDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/http"
	authRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/repo"
	authUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/usecase"
	creatorDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/http"
	creatorRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/repo"
	creatorUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	postDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/delivery/http"
	postRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/repo"
	postUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/usecase"
	userDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/http"
	userRepository "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/repo"
	userUsecase "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/usecase"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
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
	logger, err := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err = logger.Sync()
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
	authHandler := authDelivery.NewAuthHandler(authUse, zapSugar)

	userRepo := userRepository.NewUserRepo(db, zapSugar)
	userUse := userUsecase.NewUserUsecase(userRepo, zapSugar)
	userHandler := userDelivery.NewUserHandler(userUse, authUse)

	creatorRepo := creatorRepository.NewCreatorRepo(db, zapSugar)
	creatorUse := creatorUsecase.NewCreatorUsecase(creatorRepo, zapSugar)
	creatorHandler := creatorDelivery.NewCreatorHandler(creatorUse, authUse)

	attachmentRepo := attachmentRepository.NewAttachmentRepo(db, zapSugar)
	attachmentUse := attachmentUsecase.NewAttachmentUsecase(attachmentRepo, zapSugar)

	postRepo := postRepository.NewPostRepo(db, zapSugar)
	postUse := postUsecase.NewPostUsecase(postRepo, zapSugar)
	postHandler := postDelivery.NewPostHandler(postUse, authUse, attachmentUse, zapSugar)

	logMw := middleware.NewLoggerMiddleware(zapSugar)
	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CORSMiddleware)
	r.Use(logMw.LogRequest)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signUp", authHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/signIn", authHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodGet, http.MethodOptions)
	}

	user := r.PathPrefix("/user").Subrouter()
	{
		user.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/donate", userHandler.Donate).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updatePassword", userHandler.UpdatePassword).Methods(http.MethodPut, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updateData", userHandler.UpdateData).Methods(http.MethodPatch, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/homePage", userHandler.GetHomePage).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updateProfilePhoto", userHandler.UpdateProfilePhoto).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
	}

	creator := r.PathPrefix("/creator").Subrouter()
	{
		creator.HandleFunc("/list", creatorHandler.GetAllCreators).Methods(http.MethodGet, http.MethodOptions)
		creator.HandleFunc("/page/{creator-uuid}", creatorHandler.GetPage).Methods(http.MethodGet, http.MethodOptions)
		creator.HandleFunc("/aim/create", creatorHandler.CreateAim).Methods(http.MethodPost, http.MethodOptions)
	}

	post := r.PathPrefix("/post").Subrouter()
	{
		post.HandleFunc("/create", postHandler.CreatePost).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/edit/{post-uuid}", postHandler.EditPost).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/addAttach/{post-uuid}", postHandler.AddAttach).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/deleteAttach/{post-uuid}", postHandler.DeleteAttach).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/addLike", postHandler.AddLike).Methods(http.MethodPut, http.MethodOptions)
		post.HandleFunc("/removeLike", postHandler.RemoveLike).Methods(http.MethodPut, http.MethodOptions)
		post.HandleFunc("/delete/{post-uuid}", postHandler.DeletePost).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/get/{post-uuid}", postHandler.GetPost).Methods(http.MethodGet, http.MethodOptions)
	}

	http.Handle("/", r)
	srv := http.Server{Handler: r, Addr: ":8000"}
	return srv.ListenAndServe()
}
