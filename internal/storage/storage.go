package storage

import (
	"context"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/models"
	"github.com/MukizuL/GophKeeper/internal/storage/pgstorage"
	"go.uber.org/fx"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage.go -package=mockstorage

type Repository interface {
	CreateNewUser(ctx context.Context, login, passwordHash string) error
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}

func newRepository(cfg *config.Config, p *pgstorage.PGStorage) Repository {
	return p
}

func Provide() fx.Option {
	return fx.Provide(newRepository)
}
