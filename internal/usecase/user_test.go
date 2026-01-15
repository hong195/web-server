package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/evrone/go-clean-template/internal/repo/persistent"
	"github.com/evrone/go-clean-template/internal/usecase/user"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDeductBalance(t *testing.T) {
	t.Parallel()

	errUserNotFound := errors.New("user not found")
	errDB := errors.New("connection refused")

	tests := []struct {
		name      string
		userID    int64
		amount    float64
		mockSetup func(repo *MockUserRepo)
		wantErr   error
	}{
		{
			name:   "success",
			userID: 1,
			amount: 100.0,
			mockSetup: func(repo *MockUserRepo) {
				repo.EXPECT().
					DeductBalance(gomock.Any(), int64(1), 100.0).
					Return(nil)
			},
			wantErr: nil,
		},
		{
			name:      "zero amount",
			userID:    1,
			amount:    0.0,
			mockSetup: func(repo *MockUserRepo) {},
			wantErr:   user.ErrInvalidAmount,
		},
		{
			name:      "negative amount",
			userID:    1,
			amount:    -50.0,
			mockSetup: func(repo *MockUserRepo) {},
			wantErr:   user.ErrInvalidAmount,
		},
		{
			name:   "insufficient funds",
			userID: 1,
			amount: 100.0,
			mockSetup: func(repo *MockUserRepo) {
				repo.EXPECT().
					DeductBalance(gomock.Any(), int64(1), 100.0).
					Return(persistent.ErrInsufficientFunds)
			},
			wantErr: persistent.ErrInsufficientFunds,
		},
		{
			name:   "user not found",
			userID: 999,
			amount: 100.0,
			mockSetup: func(repo *MockUserRepo) {
				repo.EXPECT().
					DeductBalance(gomock.Any(), int64(999), 100.0).
					Return(errUserNotFound)
			},
			wantErr: errUserNotFound,
		},
		{
			name:   "db error",
			userID: 1,
			amount: 100.0,
			mockSetup: func(repo *MockUserRepo) {
				repo.EXPECT().
					DeductBalance(gomock.Any(), int64(1), 100.0).
					Return(errDB)
			},
			wantErr: errDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := NewMockUserRepo(ctrl)
			tt.mockSetup(repo)

			uc := user.New(repo)
			err := uc.DeductBalance(context.Background(), tt.userID, tt.amount)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
