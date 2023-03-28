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
	"log"
	"net/http"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Print(err)
		os.Exit(1)
	}
}

func run() error {
	str, err := utils.GetConnectionString()
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", str)
	if err != nil {
		return err
	}
	defer db.Close()
	db.SetMaxOpenConns(10)

	tokenGenerator := authUsecase.NewTokenator()
	encryptor, err := authUsecase.NewEncryptor()
	if err != nil {
		return err
	}

	authRepo := authRepository.NewAuthRepo(db)
	authUse := authUsecase.NewAuthUsecase(authRepo, tokenGenerator, encryptor)
	authHandler := authDelivery.NewAuthHandler(authUse)

	userRepo := userRepository.NewUserRepo(db)
	userUse := userUsecase.NewUserUsecase(userRepo)
	userHandler := userDelivery.NewUserHandler(userUse, authUse)

	creatorRepo := creatorRepository.NewCreatorRepo(db)
	creatorUse := creatorUsecase.NewCreatorUsecase(creatorRepo)
	creatorHandler := creatorDelivery.NewCreatorHandler(creatorUse)

	attachmentRepo := attachmentRepository.NewAttachmentRepo(db)
	attachmentUse := attachmentUsecase.NewAttachmentUsecase(attachmentRepo)

	postRepo := postRepository.NewPostRepo(db)
	postUse := postUsecase.NewPostUsecase(postRepo)
	postHandler := postDelivery.NewPostHandler(postUse, authUse, attachmentUse)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CORSMiddleware)
	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signUp", authHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/signIn", authHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodGet, http.MethodOptions)
	}

	user := r.PathPrefix("/user").Subrouter()
	{
		user.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/homePage", userHandler.GetHomePage).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updateProfilePhoto", userHandler.UpdateProfilePhoto).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
	}

	creator := r.PathPrefix("/creator").Subrouter()
	{
		creator.HandleFunc("/page/{creator-uuid}", creatorHandler.GetPage).Methods(http.MethodGet, http.MethodOptions)
	}

	post := r.PathPrefix("/post").Subrouter()
	{
		post.HandleFunc("/create", postHandler.CreatePost).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/addLike", postHandler.AddLike).Methods(http.MethodPut, http.MethodOptions)
		post.HandleFunc("/removeLike", postHandler.RemoveLike).Methods(http.MethodPut, http.MethodOptions)
		post.HandleFunc("/delete/{post-uuid}", postHandler.DeletePost).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		post.HandleFunc("/get/{post-uuid}", postHandler.GetPost).Methods(http.MethodGet, http.MethodOptions, http.MethodGet)
	}

	http.Handle("/", r)
	srv := http.Server{Handler: r, Addr: ":8000"}
	return srv.ListenAndServe()
}
