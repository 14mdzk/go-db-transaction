package repository

import (
	"context"

	"github.com/14mdzk/exp/internal/adapter/db_adapter"
	"github.com/14mdzk/exp/internal/user/domain"
)

type UserActivityUowInterface interface {
	WithTx(tx db_adapter.Transaction) UserActivityUowInterface
	CreateUser(ctx context.Context, user *domain.User) error
	CreateUserActivity(ctx context.Context, activity *domain.UserActivity) error
}

type UserActivityUow struct {
	userRepo     UserRepositoryInterface
	activityRepo UserActivityRepositoryInterface
}

func NewUserActivityUow(userRepo UserRepositoryInterface, activityRepo UserActivityRepositoryInterface) *UserActivityUow {
	return &UserActivityUow{userRepo: userRepo, activityRepo: activityRepo}
}

func (a *UserActivityUow) WithTx(tx db_adapter.Transaction) UserActivityUowInterface {
	return &UserActivityUow{userRepo: a.userRepo.WithTx(tx), activityRepo: a.activityRepo.WithTx(tx)}
}

func (a *UserActivityUow) CreateUser(ctx context.Context, user *domain.User) error {
	return a.userRepo.CreateUser(ctx, user)
}

func (a *UserActivityUow) CreateUserActivity(ctx context.Context, activity *domain.UserActivity) error {
	return a.activityRepo.CreateUserActivity(ctx, activity)
}
