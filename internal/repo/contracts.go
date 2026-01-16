// Package repo implements application outer layer logic. Each logic group in own file.
package repo

import (
	"context"

	"github.com/hong195/web-server/internal/entity"
)

//go:generate mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

type (
	// UserRepo -.
	UserRepo interface {
		GetByID(ctx context.Context, userID int64) (*entity.User, error)
		DeductBalance(ctx context.Context, userID int64, amount float64) error
	}

	// ItemsRepo - источник данных для items (внешний API).
	ItemsRepo interface {
		GetItems(ctx context.Context) ([]entity.Item, error)
	}
)
