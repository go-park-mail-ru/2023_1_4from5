package models

import (
	generatedCreator "github.com/go-park-mail-ru/2023_1_4from5/internal/pkg/creator/delivery/grpc/generated"
	"github.com/google/uuid"
	"html"
)

// easyjson -all ./internal/models/creator.go

type Creator struct {
	Id             uuid.UUID `json:"creator_id"`
	UserId         uuid.UUID `json:"user_id"`
	Name           string    `json:"name"`
	CoverPhoto     uuid.UUID `json:"cover_photo"`
	ProfilePhoto   uuid.UUID `json:"profile_photo"`
	FollowersCount int64     `json:"followers_count"`
	Description    string    `json:"description"`
	PostsCount     int64     `json:"posts_count"`
}

type CreatorPage struct {
	CreatorInfo   Creator        `json:"creator_info"`
	Aim           Aim            `json:"aim"`
	IsMyPage      bool           `json:"is_my_page"`
	Follows       bool           `json:"follows"`
	Posts         []Post         `json:"posts"`
	Subscriptions []Subscription `json:"subscriptions"`
}

type Aim struct {
	Creator     uuid.UUID `json:"creator_id"`
	Description string    `json:"description"`
	MoneyNeeded float32   `json:"money_needed"`
	MoneyGot    float32   `json:"money_got"`
}

type UpdateCreatorInfo struct {
	Description string    `json:"description"`
	CreatorName string    `json:"creator_name"`
	CreatorID   uuid.UUID `json:"-"`
}

func (creator *Creator) Sanitize() {
	creator.Name = html.EscapeString(creator.Name)
	creator.Description = html.EscapeString(creator.Description)
}

func (aim *Aim) Sanitize() {
	aim.Description = html.EscapeString(aim.Description)
}

func (page *CreatorPage) Sanitize() {
	page.CreatorInfo.Sanitize()
	page.Aim.Sanitize()
	for i := range page.Posts {
		page.Posts[i].Sanitize()
	}
	for i := range page.Subscriptions {
		page.Subscriptions[i].Sanitize()
	}
}

func (aim *Aim) AimToModel(aimProto *generatedCreator.Aim) error {
	creatorID, err := uuid.Parse(aimProto.Creator)
	if err != nil {
		return err
	}

	aim.Creator = creatorID
	aim.Description = aimProto.Description
	aim.MoneyNeeded = aimProto.MoneyNeeded
	aim.MoneyGot = aimProto.MoneyGot
	return nil
}
func (creator *Creator) CreatorToModel(creatorInfo *generatedCreator.Creator) error {
	creatorPhoto, err := uuid.Parse(creatorInfo.CreatorPhoto)
	if err != nil {
		return err
	}
	coverPhoto, err := uuid.Parse(creatorInfo.CoverPhoto)
	if err != nil {
		return err
	}
	creatorID, err := uuid.Parse(creatorInfo.Id)
	if err != nil {
		return err
	}
	userID, err := uuid.Parse(creatorInfo.UserID)
	if err != nil {
		return err
	}

	creator.Id = creatorID
	creator.UserId = userID
	creator.Name = creatorInfo.CreatorName
	creator.CoverPhoto = coverPhoto
	creator.ProfilePhoto = creatorPhoto
	creator.FollowersCount = creatorInfo.FollowersCount
	creator.Description = creatorInfo.Description
	creator.PostsCount = creatorInfo.PostsCount
	return nil
}
