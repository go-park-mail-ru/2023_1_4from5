package notifications

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"strconv"
)

// topic name is creator_id for users
// topic name for creator is user_id

func (na *NotificationApp) AddUserToNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	registrationTokens := []string{token.Token}

	response, err := client.SubscribeToTopic(ctx, registrationTokens, topic)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	fmt.Println(strconv.Itoa(response.SuccessCount) + " tokens were subscribed successfully")
	return nil
}

func (na *NotificationApp) RemoveUserFromNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	registrationTokens := []string{token.Token}

	response, err := client.UnsubscribeFromTopic(ctx, registrationTokens, topic)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	fmt.Println(strconv.Itoa(response.SuccessCount) + " tokens were subscribed successfully")
	return nil
}

func (na *NotificationApp) SendUserNotification(topic string, body string, title string, ctx context.Context) error {
	client, err := na.app.Messaging(ctx)
	if err != nil {
		na.logger.Error(err)
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
	}
	response, err := client.Send(ctx, message)
	if err != nil {
		na.logger.Error(err)
		return err
	}
	fmt.Println("Successfully sent message:", response)
	return nil
}
