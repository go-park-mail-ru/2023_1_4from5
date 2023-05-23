package usecase

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"go.uber.org/zap"
	"google.golang.org/api/option"
)

type NotificationApp struct {
	logger *zap.SugaredLogger
	app    *firebase.App
}

func SetupFirebase(ctx context.Context, logger *zap.SugaredLogger) *NotificationApp {
	opt := option.WithCredentialsFile("./service.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		logger.Error(err)
	}

	return &NotificationApp{
		logger: logger,
		app:    app,
	}
}

// topic name is creator_id for users
// topic name for creator is user_id

func (na *NotificationApp) AddUserToNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	registrationTokens := []string{token.Token}

	_, err = client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	return nil
}

func (na *NotificationApp) RemoveUserFromNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	registrationTokens := []string{token.Token}

	_, err = client.UnsubscribeFromTopic(ctx, registrationTokens, topic)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	return nil
}

func (na *NotificationApp) SendUserNotification(notification models.Notification, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.Photo,
		},
		Topic: notification.Topic,
	}
	_, err = client.Send(ctx, message)
	if err != nil {
		na.logger.Error(err)
		return err
	}
	return nil
}
