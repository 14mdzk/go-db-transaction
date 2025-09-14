package repository

import (
	"context"

	"github.com/14mdzk/exp/internal/adapter/db_adapter"
	"github.com/14mdzk/exp/internal/user/domain"
)

type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *domain.User) error
	WithTx(tx db_adapter.Transaction) *UserRepository
}

type UserRepository struct {
	db db_adapter.Query
}

func NewUserRepository(db db_adapter.Query) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) WithTx(tx db_adapter.Transaction) *UserRepository {
	return &UserRepository{db: tx}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	row := r.db.QueryRow(
		ctx,
		"INSERT INTO users (username, email, password, is_active) VALUES ($1, $2, $3, $4) RETURNING id",
		user.Username, user.Email, user.Password, user.IsActive,
	)

	if err := row.Scan(&user.ID); err != nil {
		return err
	}

	return nil
}
