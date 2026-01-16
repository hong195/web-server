package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/hong195/web-server/internal/entity"
	"github.com/hong195/web-server/internal/repo"
)

var (
	ErrInvalidAmount = errors.New("amount must be greater than zero")
	ErrUserNotFound  = errors.New("user not found")
)

type UseCase struct {
	repo repo.UserRepo
}

func New(r repo.UserRepo) *UseCase {
	return &UseCase{repo: r}
}

func (uc *UseCase) GetByID(ctx context.Context, userID int64) (*entity.User, error) {
	user, err := uc.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetByID: %w", err)
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
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
