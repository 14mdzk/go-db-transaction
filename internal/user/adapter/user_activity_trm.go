package adapter

import (
	"github.com/14mdzk/exp/internal/adapter/db_adapter"
	"github.com/14mdzk/exp/internal/user/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserActivityTrm(pool *pgxpool.Pool, uow repository.UserActivityUowInterface) db_adapter.Trm[repository.UserActivityUowInterface] {
	return db_adapter.NewTrm(pool, uow)
}
