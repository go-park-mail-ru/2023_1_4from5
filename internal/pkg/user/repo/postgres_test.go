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

var userID = uuid.New()
var path = uuid.New()
var password = "1234567aa"
var profileInfo = models.UpdateProfileInfo{Login: "Dasha2003!", Name: "Taktashova Daria"}

func TestUserRepo_UpdateProfilePhoto(t *testing.T) {
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
	r := NewUserRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		userID      uuid.UUID
		path        uuid.UUID
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectQuery(`UPDATE "user" SET profile_photo \= \$1 WHERE`).
					WithArgs(path, userID).WillReturnError(sql.ErrNoRows)
			},
			userID: userID,
			path:   path,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`UPDATE "user" SET profile_photo \= \$1 WHERE`).
					WithArgs(path, userID).WillReturnError(errors.New("test"))
			},
			userID:      userID,
			path:        path,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.UpdateProfilePhoto(context.Background(), test.userID, test.path)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_UpdatePassword(t *testing.T) {
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
	r := NewUserRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		userID      uuid.UUID
		password    string
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectQuery(` UPDATE "user" SET password_hash \= \$1, user_version \= user_version\+1 WHERE`).
					WithArgs(password, userID).WillReturnError(sql.ErrNoRows)
			},
			userID:   userID,
			password: password,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(` UPDATE "user" SET password_hash \= \$1, user_version \= user_version\+1 WHERE`).
					WithArgs(password, userID).WillReturnError(errors.New("test"))
			},
			userID:      userID,
			password:    password,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.UpdatePassword(context.Background(), test.userID, test.password)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepo_UpdateProfileInfo(t *testing.T) {
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
	r := NewUserRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		userID      uuid.UUID
		profile     models.UpdateProfileInfo
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				mock.ExpectQuery(`UPDATE "user" SET login \= \$1, display_name \= \$2 WHERE user_id \= \$3`).
					WithArgs(profileInfo.Login, profileInfo.Name, userID).WillReturnError(sql.ErrNoRows)
			},
			userID:  userID,
			profile: profileInfo,
		},
		{
			name: "InternalError",
			mock: func() {
				mock.ExpectQuery(`UPDATE "user" SET login \= \$1, display_name \= \$2 WHERE user_id \= \$3`).
					WithArgs(profileInfo.Login, profileInfo.Name, userID).WillReturnError(errors.New("test"))
			},
			userID:      userID,
			profile:     profileInfo,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			err := r.UpdateProfileInfo(context.Background(), test.profile, test.userID)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
