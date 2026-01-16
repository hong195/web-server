package usecase

import (
	"context"

	"github.com/hong195/web-server/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	User interface {
		GetByID(ctx context.Context, userID int64) (*entity.User, error)
		DeductBalance(ctx context.Context, userID int64, amount float64) error
	}

	Items interface {
		GetItems(ctx context.Context) ([]entity.Item, error)
	}
)
