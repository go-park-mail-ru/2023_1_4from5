package repo

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestCreatorRepo_GetUserSubscriptions(t *testing.T) {
	id := "566ece0a-a3a4-466c-8425-251147a68e90"
	sub, _ := uuid.Parse(id)
	var testUserId = uuid.New()
	var testSubscriptionId = []uuid.UUID{sub}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err.Error())
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	r := NewCreatorRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		input       uuid.UUID
		expectedRes []uuid.UUID
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"subscription_id"}).AddRow("{'" + id + "'}")
				mock.ExpectQuery(`SELECT array_agg\(subscription_id\) FROM "user_subscription" WHERE `).
					WithArgs(testUserId).WillReturnRows(rows)
			},
			input:       testUserId,
			expectedRes: testSubscriptionId,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`SELECT array_agg\(subscription_id\) FROM "user_subscription" WHERE `).
					WithArgs(testUserId).WillReturnError(errors.New("test"))
			},
			input:       testUserId,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetUserSubscriptions(context.Background(), test.input)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes[0], got[0])
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatorRepo_IsLiked(t *testing.T) {
	userId := uuid.New()
	postId := uuid.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err.Error())
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	r := NewCreatorRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		userId      uuid.UUID
		postId      uuid.UUID
		expectedRes bool
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"post_id", "user_id"}).AddRow(postId, userId)
				mock.ExpectQuery(`SELECT post_id, user_id FROM "like_post" WHERE`).
					WithArgs(postId, userId).WillReturnRows(rows)
			},
			userId:      userId,
			postId:      postId,
			expectedRes: true,
			expectedErr: nil,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`SELECT post_id, user_id FROM "like_post" WHERE`).
					WithArgs(postId, userId).WillReturnError(errors.New("test"))
			},
			userId:      userId,
			postId:      postId,
			expectedRes: false,
			expectedErr: models.InternalError,
		},
		{
			name: "Not Liked",
			mock: func() {
				mock.ExpectQuery(`SELECT post_id, user_id FROM "like_post" WHERE`).
					WithArgs(postId, userId).WillReturnError(sql.ErrNoRows)
			},
			userId:      userId,
			postId:      postId,
			expectedRes: false,
			expectedErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.IsLiked(context.Background(), test.userId, test.postId)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreatorRepo_CreatorInfo(t *testing.T) {
	//creatorPage := &models.CreatorPage{}
	creatorInfo := models.Creator{UserId: uuid.New(), Name: "testName", FollowersCount: 5, Description: "test", PostsCount: 10}
	creatorId := uuid.New()
	creatorAim := models.Aim{MoneyGot: 100, MoneyNeeded: 200, Description: "testAim", Creator: creatorId}
	creatorPageRes := models.CreatorPage{CreatorInfo: creatorInfo, Aim: creatorAim}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		t.Error(err.Error())
	}
	defer func(logger *zap.Logger) {
		err = logger.Sync()
		if err != nil {
			return
		}
	}(logger)
	zapSugar := logger.Sugar()
	r := NewCreatorRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		creatorPage *models.CreatorPage
		creatorId   uuid.UUID
		expectedRes models.CreatorPage
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "name", "cover_photo", "followers_count", "description", "posts_count",
					"aim", "money_got", "money_needed", "profile_photo"}).AddRow(creatorInfo.UserId, creatorInfo.Name, creatorInfo.CoverPhoto, creatorInfo.FollowersCount,
					creatorInfo.Description, creatorInfo.PostsCount, creatorAim.Description, creatorAim.MoneyGot, creatorAim.MoneyNeeded, creatorInfo.ProfilePhoto)
				mock.ExpectQuery(`SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)
			},
			creatorPage: &models.CreatorPage{},
			creatorId:   creatorId,
			expectedRes: creatorPageRes,
			expectedErr: nil,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE`).
					WithArgs(creatorId).WillReturnError(errors.New("test"))
			},
			creatorPage: &models.CreatorPage{},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE`).
					WithArgs(creatorId).WillReturnError(sql.ErrNoRows)
			},
			creatorPage: &models.CreatorPage{},
			creatorId:   creatorId,
			expectedErr: models.NotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.CreatorInfo(context.Background(), test.creatorPage, test.creatorId)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes, *test.creatorPage)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
