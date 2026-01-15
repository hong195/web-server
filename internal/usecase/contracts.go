// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"context"
)

//go:generate mockgen -source=contracts.go -destination=./mocks_usecase_test.go -package=usecase_test

type (
	// User -.
	User interface {
		DeductBalance(ctx context.Context, userID int64, amount float64) error
	}
)
