package persistent

import (
	"context"
	"errors"
	"fmt"

	"github.com/evrone/go-clean-template/internal/entity"
	"github.com/evrone/go-clean-template/pkg/postgres"
	"github.com/jackc/pgx/v5"
)

var ErrInsufficientFunds = errors.New("insufficient funds")

// UserRepo -.
type UserRepo struct {
	*postgres.Postgres
}

// NewUserRepo -.
func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

// GetByID -.
func (r *UserRepo) GetByID(ctx context.Context, userID int64) (*entity.User, error) {
	sql, args, err := r.Builder.
		Select("id, balance").
		From("users").
		Where("id = ?", userID).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("UserRepo - GetByID - r.Builder: %w", err)
	}

	var user entity.User
	err = r.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID, &user.Balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepo - GetByID - r.Pool.QueryRow: %w", err)
	}

	return &user, nil
}

// DeductBalance -.
func (r *UserRepo) DeductBalance(ctx context.Context, userID int64, amount float64) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - DeductBalance - r.Pool.Begin: %w", err)
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	var balance float64
	err = tx.QueryRow(ctx, "SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("UserRepo - DeductBalance - user not found: %w", err)
		}
		return fmt.Errorf("UserRepo - DeductBalance - tx.QueryRow: %w", err)
	}

	// Check sufficient funds
	if balance < amount {
		return ErrInsufficientFunds
	}

	// Update balance
	_, err = tx.Exec(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)
	if err != nil {
		return fmt.Errorf("UserRepo - DeductBalance - tx.Exec: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo - DeductBalance - tx.Commit: %w", err)
	}

	return nil
}
