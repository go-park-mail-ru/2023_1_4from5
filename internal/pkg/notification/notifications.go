package notification

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

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
