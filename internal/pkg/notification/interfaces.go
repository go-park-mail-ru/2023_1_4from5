package notification

import (
	"context"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks/notification_mock.go -package=mock

type NotificationApp interface {
	SendUserNotification(notification models.Notification, ctx context.Context) error
	RemoveUserFromNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error
	AddUserToNotificationTopic(topic string, token models.NotificationToken, ctx context.Context) error
}
