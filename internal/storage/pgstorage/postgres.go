package pgstorage

import (
	"context"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

type PGStorage struct {
	conn *pgxpool.Pool
	cfg  *config.Config
}

func newPGStorage(cfg *config.Config) *PGStorage {
	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		panic(err)
	}

	return &PGStorage{
		conn: dbpool,
		cfg:  cfg,
	}
}

func Provide() fx.Option {
	return fx.Provide(newPGStorage)
}
