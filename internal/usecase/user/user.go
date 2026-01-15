package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/evrone/go-clean-template/internal/repo"
)

var ErrInvalidAmount = errors.New("amount must be greater than zero")

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) DeductBalance(ctx context.Context, userID int64, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	err := uc.repo.DeductBalance(ctx, userID, amount)
	if err != nil {
		return fmt.Errorf("UserUseCase - DeductBalance: %w", err)
	}

	return nil
}
