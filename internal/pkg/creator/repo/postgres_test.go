package repo

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestCreatorRepo_GetUserSubscriptions(t *testing.T) {
	sub, _ := uuid.Parse("566ece0a-a3a4-466c-8425-251147a68e90")
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
				rows := sqlmock.NewRows([]string{"subscription_id"}).AddRow("{'566ece0a-a3a4-466c-8425-251147a68e90'}")
				mock.ExpectQuery(`SELECT array_agg\(subscription_id\) FROM "user_subscription" WHERE `).
					WithArgs(testUserId).WillReturnRows(rows)
			},
			input:       testUserId,
			expectedRes: testSubscriptionId,
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
