package notifications

import (
	"context"
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

//func SendToToken(app *firebase.App) {
//	ctx := context.Background()
//	client, err := app.Messaging(ctx)
//	if err != nil {
//		log.Fatalf("error getting Messaging client: %v\n", err)
//	}
//
//	registrationToken := "device token we got from frontend"
//
//	message := &messaging.Message{
//		Notification: &messaging.Notification{
//			Title: "Notification Test",
//			Body:  "Hello!!",
//		},
//		Token: registrationToken,
//	}
//
//	response, err := client.Send(ctx, message)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println("Successfully sent message:", response)
//}
