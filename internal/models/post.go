package models

// easyjson -all ./internal/models/post.go

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	Id          uuid.UUID `json:"id"`
	Creator     uuid.UUID `json:"creator"`
	Creation    time.Time `json:"creation_date"`
	Title       string    `json:"title"`
	Text        string    `json:"text"`
	IsAvailable bool      `json:"is_available"`
}

//После signIn возвращаем на запрос о пользователе
//user_name
//ленту
//is_creator:true/false
//if is_creator:true, то creator_id
//
//Страница автора
//Получаем на бэк
//user_id
//creator_id
//
//Отдаём на фронт
//myPage:true/false (проверяем вдруг автор зашёл на свою страницу)
//posts
//subscriptions
