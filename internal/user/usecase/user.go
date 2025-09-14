package usecase

import (
	"context"

	"github.com/14mdzk/exp/internal/adapter/db_adapter"
	"github.com/14mdzk/exp/internal/user/domain"
	"github.com/14mdzk/exp/internal/user/dto"
	"github.com/14mdzk/exp/internal/user/repository"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userTrm          db_adapter.Trm[repository.UserActivityUowInterface]
	userActivityRepo repository.UserActivityRepositoryInterface
}

func NewUserUsecase(userTrm db_adapter.Trm[repository.UserActivityUowInterface], userActivityRepo repository.UserActivityRepositoryInterface) *userUsecase {
	return &userUsecase{userTrm, userActivityRepo}
}

func (u *userUsecase) CreateUser(ctx context.Context, user *dto.CreateUser) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	hashedStr := string(hashed)
	newUser := &domain.User{
		Username: user.Username,
		Email:    user.Email,
		Password: &hashedStr,
	}

	err = u.userTrm.InTx(ctx, func(tx repository.UserActivityUowInterface) error {
		if err := tx.CreateUser(ctx, newUser); err != nil {
			u.userActivityRepo.CreateUserActivity(ctx, domain.NewUserActivity(nil, domain.UserActivityNameCreate, err.Error()))
			return err
		}

		if err := tx.CreateUserActivity(ctx, domain.NewUserActivity(nil, domain.UserActivityNameCreate, "create user with id "+newUser.ID.String())); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := u.userActivityRepo.CreateUserActivity(ctx, domain.NewUserActivity(nil, domain.UserActivityNameCreate, "create user with id "+newUser.ID.String()+" success")); err != nil {
		return err
	}

	return nil
}
