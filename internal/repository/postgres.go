package repository

import (
	"bot/internal"

	"github.com/jackc/pgx/v4/pgxpool"
)

type dataRepo struct {
	pool *pgxpool.Pool
}

func NewDataRepo(pool *pgxpool.Pool) internal.Repository {
	return &dataRepo{
		pool: pool,
	}
}