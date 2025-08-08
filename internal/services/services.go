package services

import (
	jwtService "github.com/MukizuL/GophKeeper/internal/jwt"
	"github.com/MukizuL/GophKeeper/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Services struct {
	storage    storage.Repository
	jwtService jwtService.ServiceI
	logger     *zap.Logger
}

func newServices(storage storage.Repository, jwtService jwtService.ServiceI, logger *zap.Logger) *Services {
	return &Services{
		storage:    storage,
		jwtService: jwtService,
		logger:     logger,
	}
}

func Provide() fx.Option {
	return fx.Provide(newServices)
}
