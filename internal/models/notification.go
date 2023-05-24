package models

import "github.com/google/uuid"

// easyjson -all ./internal/models/notification.go
const PhotoURL = "https://sub-me.ru/images/user/"

type NotificationToken struct {
	Token string `json:"notification_token"`
}

//easyjson:skip
type Notification struct {
	Topic string
	Title string
	Body  string
	Photo string
}

//easyjson:skip
type NotificationCreatorInfo struct {
	Name  string
	Photo uuid.UUID
}

//easyjson:skip
type NotificationSubInfo struct {
	SubscriptionName string
	CreatorID        uuid.UUID
}
