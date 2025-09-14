package repository

import (
	"context"

	"github.com/14mdzk/exp/internal/adapter/db_adapter"
	"github.com/14mdzk/exp/internal/user/domain"
)

type UserActivityRepositoryInterface interface {
	CreateUserActivity(ctx context.Context, activity *domain.UserActivity) error
	WithTx(tx db_adapter.Transaction) *UserActivityRepository
}

type UserActivityRepository struct {
	db db_adapter.Query
}

func NewUserActivityRepository(db db_adapter.Query) *UserActivityRepository {
	return &UserActivityRepository{db}
}

func (r *UserActivityRepository) WithTx(tx db_adapter.Transaction) *UserActivityRepository {
	return &UserActivityRepository{tx}
}

func (r *UserActivityRepository) CreateUserActivity(ctx context.Context, activity *domain.UserActivity) error {
	row := r.db.QueryRow(
		ctx,
		"INSERT INTO activities (object, object_id, name, description) VALUES ($1, $2, $3, $4) RETURNING id",
		activity.Object, activity.ObjectID, string(activity.Name), activity.Description,
	)

	if err := row.Scan(&activity.ID); err != nil {
		return err
	}

	return nil
}
