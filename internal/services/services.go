package services

import (
	"context"

	jwtService "github.com/MukizuL/GophKeeper/internal/jwt"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/MukizuL/GophKeeper/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:generate mockgen -source=services.go -destination=mocks/services.go -package=mockservices

type ServicesI interface {
	CreateNewUser(ctx context.Context, login, password string) error
	Login(ctx context.Context, login, password string) (string, []byte, error)
	CreatePassword(ctx context.Context, data []byte) error
	CreateBank(ctx context.Context, data []byte) error
	CreateTextual(ctx context.Context, data []byte) error
	CreateData(ctx context.Context, token string, stream pb.Gophkeeper_CreateDataServer) error
	GetPasswords(ctx context.Context) ([][]byte, error)
	GetBank(ctx context.Context) ([][]byte, error)
	GetText(ctx context.Context) ([][]byte, error)
	GetData(ctx context.Context) ([]*pb.File, error)
	Download(ctx context.Context, token, id string, stream pb.Gophkeeper_DownloadServer) error
}

type Services struct {
	storage    storage.Repository
	jwtService jwtService.ServiceI
	logger     *zap.Logger
}

func newServices(storage storage.Repository, jwtService jwtService.ServiceI, logger *zap.Logger) ServicesI {
	return &Services{
		storage:    storage,
		jwtService: jwtService,
		logger:     logger,
	}
}

func Provide() fx.Option {
	return fx.Provide(newServices)
}
