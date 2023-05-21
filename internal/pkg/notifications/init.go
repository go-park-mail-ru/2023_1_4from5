package notifications

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type NotificationApp struct {
	logger *zap.SugaredLogger
	app    *firebase.App
}

func SetupFirebase(ctx context.Context, logger *zap.SugaredLogger) *NotificationApp {
	opt := option.WithCredentialsFile("/home/dasha/Downloads/service.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Error(err)
	}

	return &NotificationApp{
		logger: logger,
		app:    app,
	}
}
