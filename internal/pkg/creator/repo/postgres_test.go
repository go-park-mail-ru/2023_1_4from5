package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var userId = uuid.New()
var postId = uuid.New()
var attachsIDs = []uuid.UUID{uuid.New(), uuid.New()}
var creatorId = uuid.New()
var subsIDs = []uuid.UUID{uuid.New(), uuid.New()}
var testSubscriptionId = []uuid.UUID{subsIDs[0]}
var id = subsIDs[0].String()
var attachTypes = []string{"test1", "test2"}
var attachments = []models.Attachment{{Id: attachsIDs[0], Type: attachTypes[0]}, {Id: attachsIDs[1], Type: attachTypes[1]}}
var subs = []models.Subscription{{Id: subsIDs[0], Creator: creatorId, MonthCost: int64(100), Title: "test", Description: "TEST"}, {Id: subsIDs[1], Creator: creatorId, MonthCost: 100}}
var posts = []models.Post{{Id: uuid.New(), Creator: creatorId, LikesCount: int64(4), Title: "test", Text: "TEST", Attachments: attachments, Subscriptions: subs}, {Id: uuid.New(), Creator: creatorId, LikesCount: 15, Title: "test1", Text: "TEST1", Attachments: attachments, Subscriptions: subs}}
var creatorInfo = models.Creator{Id: creatorId, UserId: uuid.New(), Name: "testName", FollowersCount: int64(5), Description: "test", PostsCount: 10}
var creatorAim = models.Aim{MoneyGot: int64(100), MoneyNeeded: int64(200), Description: "testAim", Creator: creatorId}
var creatorPageRes = models.CreatorPage{CreatorInfo: creatorInfo, Aim: creatorAim}
var creatorPageRes2 = models.CreatorPage{CreatorInfo: creatorInfo, Aim: creatorAim, Posts: posts, Subscriptions: subs}
var creators = []models.Creator{creatorInfo, creatorInfo}

func TestCreatorRepo_GetUserSubscriptions(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
					WithArgs(userId).WillReturnRows(rows)
			},
			input:       userId,
			expectedRes: testSubscriptionId,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`SELECT array_agg\(subscription_id\) FROM "user_subscription" WHERE `).
					WithArgs(userId).WillReturnError(errors.New("test"))
			},
			input:       userId,
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

func TestCreatorRepo_CreatorPosts(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		creatorId   uuid.UUID
		expectedRes []models.Post
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"post_id", "creation_date", "title", "post_text", "likes_count", "attachment_id", "attachment_type", "subscription_id"})

				rows = rows.AddRow(posts[0].Id, posts[0].Creation, posts[0].Title, posts[0].Text, posts[0].LikesCount, fmt.Sprintf("{'%s','%s'}", attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s}", attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))
				rows = rows.AddRow(posts[1].Id, posts[1].Creation, posts[1].Title, posts[1].Text, posts[1].LikesCount, fmt.Sprintf("{'%s','%s'}", attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s}", attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))

				mock.ExpectQuery(`SELECT "post"\.post_id, creation_date, title, post_text, likes_count, array_agg\(attachment_id\), array_agg\(attachment_type\), array_agg\(DISTINCT subscription_id\) FROM "post" LEFT JOIN "attachment" a on "post"\.post_id \= a\.post_id LEFT JOIN "post_subscription" ps on "post"\.post_id \= ps\.post_id WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)
				for i := 0; i < 4; i++ {
					rows = sqlmock.NewRows([]string{"creator_id", "month_cost", "title", "description"}).AddRow(subs[i%2].Creator, subs[i%2].MonthCost, subs[i%2].Title, subs[i%2].Description)
					mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[i%2]).WillReturnRows(rows)
				}
			},
			creatorId:   creatorId,
			expectedRes: posts,
		},
		{
			name: "Internal Error for Get Posts",
			mock: func() {
				mock.ExpectQuery(`SELECT "post"\.post_id, creation_date, title, post_text, likes_count, array_agg\(attachment_id\), array_agg\(attachment_type\), array_agg\(DISTINCT subscription_id\) FROM "post" LEFT JOIN "attachment" a on "post"\.post_id \= a\.post_id LEFT JOIN "post_subscription" ps on "post"\.post_id \= ps\.post_id WHERE`).
					WithArgs(creatorId).WillReturnError(errors.New("test"))
			},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
		{
			name: "Internal Error in GetSubsById",
			mock: func() {
				rows := sqlmock.NewRows([]string{"post_id", "creation_date", "title", "post_text", "likes_count", "attachment_id", "attachment_type", "subscription_id"})

				rows = rows.AddRow(posts[0].Id, posts[0].Creation, posts[0].Title, posts[0].Text, posts[0].LikesCount, fmt.Sprintf("{'%s','%s','%s','%s'}", attachsIDs[0], attachsIDs[1], attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s,%s,%s}", attachTypes[0], attachTypes[1], attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))
				rows = rows.AddRow(posts[1].Id, posts[1].Creation, posts[1].Title, posts[1].Text, posts[1].LikesCount, fmt.Sprintf("{'%s','%s','%s','%s'}", attachsIDs[0], attachsIDs[1], attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s,%s,%s}", attachTypes[0], attachTypes[1], attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))

				mock.ExpectQuery(`SELECT "post"\.post_id, creation_date, title, post_text, likes_count, array_agg\(attachment_id\), array_agg\(attachment_type\), array_agg\(DISTINCT subscription_id\) FROM "post" LEFT JOIN "attachment" a on "post"\.post_id \= a\.post_id LEFT JOIN "post_subscription" ps on "post"\.post_id \= ps\.post_id WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)

				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[0]).WillReturnError(models.InternalError)

			},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
		{
			name: "Internal Error wrong data type",
			mock: func() {
				rows := sqlmock.NewRows([]string{"post_id", "creation_date", "title", "post_text", "likes_count", "attachment_id", "attachment_type", "subscription_id"})

				rows = rows.AddRow(posts[0].Id, posts[0].Creation, posts[0].Title, posts[0].Text, false, fmt.Sprintf("{'%s','%s','%s','%s'}", attachsIDs[0], attachsIDs[1], attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s,%s,%s}", attachTypes[0], attachTypes[1], attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))
				rows = rows.AddRow(posts[1].Id, posts[0].Creation, posts[0].Title, posts[0].Text, posts[1].LikesCount, fmt.Sprintf("{'%s','%s','%s','%s'}", attachsIDs[0], attachsIDs[1], attachsIDs[0], attachsIDs[1]), fmt.Sprintf("{%s,%s,%s,%s}", attachTypes[0], attachTypes[1], attachTypes[0], attachTypes[1]), fmt.Sprintf("{'%s','%s'}", subsIDs[0], subsIDs[1]))

				mock.ExpectQuery(`SELECT "post"\.post_id, creation_date, title, post_text, likes_count, array_agg\(attachment_id\), array_agg\(attachment_type\), array_agg\(DISTINCT subscription_id\) FROM "post" LEFT JOIN "attachment" a on "post"\.post_id \= a\.post_id LEFT JOIN "post_subscription" ps on "post"\.post_id \= ps\.post_id WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)

			},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CreatorPosts(context.Background(), test.creatorId)
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

func TestCreatorRepo_GetSubsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		subIDs      []uuid.UUID
		expectedRes []models.Subscription
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"creator_id", "month_cost", "title", "description"}).AddRow(subs[0].Creator, subs[0].MonthCost, subs[0].Title, subs[0].Description)
				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[0]).WillReturnRows(rows)
				rows = sqlmock.NewRows([]string{"creator_id", "month_cost", "title", "description"}).AddRow(subs[1].Creator, subs[1].MonthCost, subs[1].Title, subs[1].Description)
				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[1]).WillReturnRows(rows)
			},
			subIDs:      subsIDs,
			expectedRes: subs,
		},
		{
			name: "Internal Error",
			mock: func() {
				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[0]).WillReturnError(errors.New("test"))
			},
			subIDs:      subsIDs,
			expectedRes: nil,
			expectedErr: models.InternalError,
		},
		{
			name: "One of IDs is invalid",
			mock: func() {
				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[0]).WillReturnError(sql.ErrNoRows)
				rows := sqlmock.NewRows([]string{"creator_id", "month_cost", "title", "description"}).AddRow(subs[1].Creator, subs[1].MonthCost, subs[1].Title, subs[1].Description)
				mock.ExpectQuery(`SELECT creator_id, month_cost, title, description FROM "subscription" WHERE`).WithArgs(subsIDs[1]).WillReturnRows(rows)
			},
			subIDs:      subsIDs,
			expectedRes: subs[1:],
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetSubsByID(context.Background(), test.subIDs...)
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

func TestCreatorRepo_IsLiked(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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

			err = r.CreatorInfo(context.Background(), test.creatorPage, test.creatorId)
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

func TestCreatorRepo_GetCreatorSubs(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		creatorId   uuid.UUID
		expectedRes []models.Subscription
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"subscription_id", "month_cost", "title", "description", "is_available"})
				for _, sub := range subs {
					rows = rows.AddRow(sub.Id, sub.MonthCost, sub.Title, sub.Description, true)
				}
				mock.ExpectQuery(`SELECT subscription_id, month_cost, title, description, is_available FROM "subscription" WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)
			},
			creatorId:   creatorId,
			expectedRes: subs,
			expectedErr: nil,
		},
		{
			name: "Internal Error in GetCreatorSubs",
			mock: func() {
				rows := sqlmock.NewRows([]string{"subscription_id", "month_cost", "title", "description"})
				for _, sub := range subs {
					rows = rows.AddRow(sub.Id, sub.MonthCost, sub.Title, sub.Description)
				}
				mock.ExpectQuery(`SELECT subscription_id, month_cost, title, description, is_available FROM "subscription" WHERE`).
					WithArgs(creatorId).WillReturnError(errors.New("test"))
			},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
		{
			name: "Internal Error in data types",
			mock: func() {
				rows := sqlmock.NewRows([]string{"subscription_id", "month_cost", "title", "description", "is_available"})
				for _, sub := range subs {
					rows = rows.AddRow(sub.Id, sub.Title, sub.Title, sub.Description, true)
				}
				mock.ExpectQuery(`SELECT subscription_id, month_cost, title, description, is_available FROM "subscription" WHERE`).
					WithArgs(creatorId).WillReturnRows(rows)
			},
			creatorId:   creatorId,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetCreatorSubs(context.Background(), test.creatorId)
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

func TestCreatorRepo_CreateAim(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		aim         models.Aim
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{})
				mock.ExpectQuery(`UPDATE creator SET aim \= \$1, money_got \= \$2, money_needed \= \$3 WHERE creator_id \= \$4; `).
					WithArgs(creatorAim.Description, creatorAim.MoneyGot, creatorAim.MoneyNeeded, creatorAim.Creator).WillReturnRows(rows)
			},
			aim:         creatorAim,
			expectedErr: nil,
		},
		{
			name: "Internal Error",
			mock: func() {
				//rows := sqlmock.NewRows([]string{})
				mock.ExpectQuery(`UPDATE creator SET aim \= \$1, money_got \= \$2, money_needed \= \$3 WHERE creator_id \= \$4; `).
					WithArgs(creatorAim.Description, creatorAim.MoneyGot, creatorAim.MoneyNeeded, creatorAim.Creator).WillReturnError(errors.New("test"))
			},
			aim:         creatorAim,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err = r.CreateAim(context.Background(), test.aim)
			assert.Equal(t, err, test.expectedErr)
		})
	}
}

func TestCreatorRepo_GetAllCreators(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		expectedRes []models.Creator
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"creator_id", "user_id", "name", "cover_photo", "followers_count", "description", "posts_count", "profile_photo"})
				rows = rows.AddRow(creatorInfo.Id, creatorInfo.UserId, creatorInfo.Name, creatorInfo.CoverPhoto, creatorInfo.FollowersCount, creatorInfo.Description, creatorInfo.PostsCount, creatorInfo.ProfilePhoto)
				rows = rows.AddRow(creatorInfo.Id, creatorInfo.UserId, creatorInfo.Name, creatorInfo.CoverPhoto, creatorInfo.FollowersCount, creatorInfo.Description, creatorInfo.PostsCount, creatorInfo.ProfilePhoto)
				mock.ExpectQuery(`SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM "creator"`).
					WithArgs().WillReturnRows(rows)
			},
			expectedRes: creators,
			expectedErr: nil,
		},
		{
			name: "Internal Error",
			mock: func() {
				mock.ExpectQuery(`SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM "creator"`).
					WithArgs().WillReturnError(errors.New("test"))
			},
			expectedErr: models.InternalError,
		},
		{
			name: "Internal Error wrong data type",
			mock: func() {
				rows := sqlmock.NewRows([]string{"creator_id", "user_id", "name", "cover_photo", "followers_count", "description", "posts_count", "profile_photo"})
				rows = rows.AddRow(creatorInfo.Id, 11, creatorInfo.Name, creatorInfo.CoverPhoto, creatorInfo.FollowersCount, creatorInfo.Description, creatorInfo.PostsCount, creatorInfo.ProfilePhoto)
				rows = rows.AddRow(creatorInfo.Id, creatorInfo.UserId, creatorInfo.Name, creatorInfo.CoverPhoto, creatorInfo.FollowersCount, creatorInfo.Description, creatorInfo.PostsCount, creatorInfo.ProfilePhoto)
				mock.ExpectQuery(`SELECT creator_id, user_id, name, cover_photo, followers_count, description, posts_count, profile_photo FROM "creator"`).
					WithArgs().WillReturnRows(rows)
			},
			expectedRes: creators,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetAllCreators(context.Background())
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

func TestCreatorRepo_GetPage(t *testing.T) {
	creatorPageRes2.Posts[0].IsAvailable = true
	creatorPageRes2.Posts[0].IsLiked = true
	creatorPageRes2.Posts[1].IsAvailable = true
	creatorPageRes2.Posts[1].IsLiked = true
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	logger := zap.NewNop()
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
		creatorID   uuid.UUID
		userID      uuid.UUID
		expectedRes models.CreatorPage
		expectedErr error
	}{
		{
			name: "Wrong Data no such author",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE`).
					WithArgs(creatorId).WillReturnError(sql.ErrNoRows)

			},
			creatorID:   creatorId,
			userID:      userId,
			expectedErr: models.WrongData,
		},
		{
			name: "Internal Error in CreatorInfo",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, name, cover_photo, followers_count, description, posts_count, aim, money_got, money_needed, profile_photo FROM "creator" WHERE`).
					WithArgs(creatorId).WillReturnError(errors.New("test"))

			},
			creatorID:   creatorId,
			userID:      userId,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.GetPage(context.Background(), test.userID, test.creatorID)
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
