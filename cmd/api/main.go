package main

import (
	"context"
	"fmt"

	"github.com/14mdzk/exp/internal/platform/config"
	"github.com/14mdzk/exp/internal/platform/db"
	"github.com/14mdzk/exp/internal/user/adapter"
	"github.com/14mdzk/exp/internal/user/dto"
	"github.com/14mdzk/exp/internal/user/repository"
	"github.com/14mdzk/exp/internal/user/usecase"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	db, err := db.NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	activityRepo := repository.NewUserActivityRepository(db)

	userActivityUow := repository.NewUserActivityUow(userRepo, activityRepo)
	userActivityTrm := adapter.NewUserActivityTrm(db.Pool, userActivityUow)

	userUsecase := usecase.NewUserUsecase(userActivityTrm, activityRepo)

	if err := runMigration(&cfg.Database); err != nil {
		fmt.Printf("migrate: %v\n", err)
	}

	fmt.Println(
		userUsecase.CreateUser(context.Background(), &dto.CreateUser{
			Username: "test_user",
			Email:    "test_user@example.com",
			Password: "test",
		}),
	)
}

func runMigration(dbConfig *config.DatabaseConfig) error {
	m, err := migrate.New(
		"file://migrations",
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name, dbConfig.SSLMode),
	)

	if err != nil {
		return err
	}

	v, d, err := m.Version()
	if err != nil {
		return err
	}

	if v != 0 && d {
		if err := m.Down(); err != nil {
			return err
		}
	}

	if err := m.Up(); err != nil {
		return err
	}

	return nil
}
