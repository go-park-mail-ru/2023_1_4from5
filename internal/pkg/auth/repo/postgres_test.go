package repo

import (
	"context"
	"database/sql"
	"fmt"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2023_1_4from5/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

var user = models.User{Id: uuid.New(), Login: "testlogin", PasswordHash: "testpwd", UserVersion: int64(2), Name: "TESTNAME", ProfilePhoto: uuid.New()}

func TestAuthRepo_CheckUser(t *testing.T) {
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
	r := NewAuthRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		input       models.User
		expectedRes models.User
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "password_hash", "user_version"}).AddRow(user.Id, user.PasswordHash, 2)
				mock.ExpectQuery(`SELECT user_id, password_hash, user_version FROM "user" WHERE`).
					WithArgs(user.Login).WillReturnRows(rows)
			},
			input:       user,
			expectedRes: user,
		},
		{
			name: "WrongPassword",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id", "password_hash", "user_version"}).AddRow(user.Id, user.PasswordHash+"1", 2)
				mock.ExpectQuery(`SELECT user_id, password_hash, user_version FROM "user" WHERE`).
					WithArgs(user.Login).WillReturnRows(rows)
			},
			input:       user,
			expectedErr: models.WrongPassword,
		},
		{
			name: "NotFound",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, password_hash, user_version FROM "user" WHERE`).
					WithArgs(user.Login).WillReturnError(sql.ErrNoRows)
			},
			input:       user,
			expectedRes: user,
			expectedErr: models.NotFound,
		},
		{
			name: "InternalErr",
			mock: func() {
				mock.ExpectQuery(`SELECT user_id, password_hash, user_version FROM "user" WHERE`).
					WithArgs(user.Login).WillReturnError(fmt.Errorf("test err"))
			},
			input:       user,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CheckUser(context.Background(), test.input)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes.Login, got.Login)
				assert.Equal(t, test.expectedRes.PasswordHash, got.PasswordHash)
				assert.Equal(t, test.expectedRes.UserVersion, got.UserVersion)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAuthRepo_CreateUser(t *testing.T) {
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
	r := NewAuthRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		input       models.User
		expectedRes models.User
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_id"}).AddRow(user.Id)
				mock.ExpectQuery(`INSERT INTO "user" \(user_id, login, display_name, profile_photo, password_hash\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING user_id;`).WithArgs(user.Id, user.Login, user.Name, user.ProfilePhoto, user.PasswordHash).WillReturnRows(rows)
			},
			input:       user,
			expectedRes: user,
		},
		{
			name: "InternalErr",
			mock: func() {
				mock.ExpectQuery(`INSERT INTO "user" \(user_id, login, display_name, profile_photo, password_hash\) VALUES\(\$1, \$2, \$3, \$4, \$5\) RETURNING user_id;`).WithArgs(user.Id, user.Login, user.Name, user.ProfilePhoto, user.PasswordHash).WillReturnError(sql.ErrNoRows)
			},
			input:       user,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CreateUser(context.Background(), test.input)
			if test.expectedErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedRes.Login, got.Login)
				assert.Equal(t, test.expectedRes.PasswordHash, got.PasswordHash)
				assert.Equal(t, test.expectedRes.UserVersion, got.UserVersion)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

var user2 = models.AccessDetails{Id: uuid.New(), Login: "testlogin", UserVersion: int64(2)}

func TestAuthRepo_CheckUserVersion(t *testing.T) {
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
	r := NewAuthRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		input       models.AccessDetails
		expectedRes int64
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_version"}).AddRow(user.UserVersion)
				mock.ExpectQuery(`SELECT user_version FROM "user" WHERE user_id \= \$1`).WithArgs(user2.Id).WillReturnRows(rows)
			},
			input:       user2,
			expectedRes: 2,
		},
		{
			name: "InternalErr",
			mock: func() {
				mock.ExpectQuery(`SELECT user_version FROM "user" WHERE user_id \= \$1`).WithArgs(user2.Id).WillReturnError(sql.ErrNoRows)
			},
			input:       user2,
			expectedErr: models.InternalError,
		},
		{
			name: "Unauthorized",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_version"}).AddRow(1)
				mock.ExpectQuery(`SELECT user_version FROM "user" WHERE user_id \= \$1`).WithArgs(user2.Id).WillReturnRows(rows)
			},
			input:       user2,
			expectedErr: models.Unauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.CheckUserVersion(context.Background(), test.input)
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

func TestAuthRepo_IncUserVersion(t *testing.T) {
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
	r := NewAuthRepo(db, zapSugar)

	tests := []struct {
		name        string
		mock        func()
		input       uuid.UUID
		expectedRes int64
		expectedErr error
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"user_version"}).AddRow(user.UserVersion)
				mock.ExpectQuery(`UPDATE "user" SET user_version \= user_version \+ 1 WHERE user_id\=\$1 RETURNING user_version;`).WithArgs(user2.Id).WillReturnRows(rows)
			},
			input:       user2.Id,
			expectedRes: int64(2),
		},
		{
			name: "InternalErr",
			mock: func() {
				mock.ExpectQuery(`UPDATE "user" SET user_version \= user_version \+ 1 WHERE user_id\=\$1 RETURNING user_version;`).WithArgs(user2.Id).WillReturnError(sql.ErrNoRows)
			},
			input:       user2.Id,
			expectedErr: models.InternalError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()

			got, err := r.IncUserVersion(context.Background(), test.input)
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
