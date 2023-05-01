package main

import (
	"database/sql"
	generatedAuth "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/grpc/generated"
	authDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/auth/delivery/http"
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	creatorDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/http"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/middleware"
	postDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/post/delivery/http"
	subscriptionDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/subscription/delivery/http"
	generatedUser "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/grpc/generated"
	userDelivery "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/user/delivery/http"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/utils"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	logger, err := utils.FileLogger("/var/log/main_app.log")
	if err != nil {
		return err
	}

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

	if err != nil {
		return err
	}

	authConn, err := grpc.Dial(
		"auth:8010",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("cant connect to session grpc")
	}

	userConn, err := grpc.Dial(
		"user:8020",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("cant connect to session grpc")
	}

	creatorConn, err := grpc.Dial(
		"creator:8030",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	if err != nil {
		log.Fatalf("cant connect to session grpc")
	}

	authClient := generatedAuth.NewAuthServiceClient(authConn)
	userClient := generatedUser.NewUserServiceClient(userConn)
	creatorClient := generatedCreator.NewCreatorServiceClient(creatorConn)

	authHandler := authDelivery.NewAuthHandler(authClient, zapSugar)
	userHandler := userDelivery.NewUserHandler(userClient, authClient, zapSugar)
	creatorHandler := creatorDelivery.NewCreatorHandler(creatorClient, authClient, zapSugar)
	postHandler := postDelivery.NewPostHandler(authClient, creatorClient, zapSugar)
	subscriptionHandler := subscriptionDelivery.NewSubscriptionHandler(authClient, creatorClient, userClient, zapSugar)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.Use(middleware.CORSMiddleware)

	logMw := middleware.NewLoggerMiddleware(zapSugar)
	r.Use(logMw.LogRequest)

	//metricsMw := middleware.NewMetricsMiddleware()
	//metricsMw.Register(middleware.ServiceMainName)
	//r.PathPrefix("/api/metrics").Handler(promhttp.Handler())
	//r.Use(metricsMw.LogMetrics)

	auth := r.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/signUp", authHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/signIn", authHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
		auth.HandleFunc("/logout", authHandler.Logout).Methods(http.MethodPut, http.MethodOptions)
	}

	user := r.PathPrefix("/user").Subrouter()
	{
		user.HandleFunc("/profile", userHandler.GetProfile).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/donate", userHandler.Donate).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updatePassword", userHandler.UpdatePassword).Methods(http.MethodPut, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updateData", userHandler.UpdateData).Methods(http.MethodPut, http.MethodGet, http.MethodOptions)
		user.HandleFunc("/feed", creatorHandler.GetFeed).Methods(http.MethodGet, http.MethodOptions)
		user.HandleFunc("/updateProfilePhoto", userHandler.UpdateProfilePhoto).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
		user.HandleFunc("/deleteProfilePhoto/{image-uuid}", userHandler.DeleteProfilePhoto).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		user.HandleFunc("/becameCreator", userHandler.BecomeCreator).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
		user.HandleFunc("/follow/{creator-uuid}", userHandler.Follow).Methods(http.MethodPost, http.MethodOptions)
		user.HandleFunc("/unfollow/{creator-uuid}", userHandler.Unfollow).Methods(http.MethodPut, http.MethodOptions)
		user.HandleFunc("/subscribe/{sub-uuid}", userHandler.Subscribe).Methods(http.MethodPost, http.MethodOptions, http.MethodGet)
		user.HandleFunc("/subscriptions", userHandler.UserSubscriptions).Methods(http.MethodOptions, http.MethodGet)
	}

	creator := r.PathPrefix("/creator").Subrouter()
	{
		creator.HandleFunc("/list", creatorHandler.GetAllCreators).Methods(http.MethodGet, http.MethodOptions)
		creator.HandleFunc("/search/{keyword}", creatorHandler.FindCreator).Methods(http.MethodGet, http.MethodOptions)
		creator.HandleFunc("/page/{creator-uuid}", creatorHandler.GetPage).Methods(http.MethodGet, http.MethodOptions)
		creator.HandleFunc("/aim/create", creatorHandler.CreateAim).Methods(http.MethodPost, http.MethodOptions)
		creator.HandleFunc("/updateData", creatorHandler.UpdateCreatorData).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
		creator.HandleFunc("/updateProfilePhoto", creatorHandler.UpdateProfilePhoto).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
		creator.HandleFunc("/deleteProfilePhoto/{image-uuid}", creatorHandler.DeleteProfilePhoto).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		creator.HandleFunc("/deleteCoverPhoto/{image-uuid}", creatorHandler.DeleteCoverPhoto).Methods(http.MethodDelete, http.MethodOptions, http.MethodGet)
		creator.HandleFunc("/updateCoverPhoto", creatorHandler.UpdateCoverPhoto).Methods(http.MethodPut, http.MethodOptions, http.MethodGet)
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

	subscription := r.PathPrefix("/subscription").Subrouter()
	{
		subscription.HandleFunc("/create", subscriptionHandler.CreateSubscription).Methods(http.MethodPost, http.MethodGet, http.MethodOptions)
		subscription.HandleFunc("/edit/{sub-uuid}", subscriptionHandler.EditSubscription).Methods(http.MethodPut, http.MethodGet, http.MethodOptions)
		subscription.HandleFunc("/delete/{sub-uuid}", subscriptionHandler.DeleteSubscription).Methods(http.MethodDelete, http.MethodGet, http.MethodOptions)
	}

	http.Handle("/", r)
	srv := http.Server{Handler: r, Addr: ":8000"}
	return srv.ListenAndServe()
}
