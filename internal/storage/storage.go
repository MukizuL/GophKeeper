package storage

import (
	"context"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/dto"
	"github.com/MukizuL/GophKeeper/internal/models"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/MukizuL/GophKeeper/internal/storage/pgstorage"
	"go.uber.org/fx"
)

//go:generate mockgen -source=storage.go -destination=mocks/storage.go -package=mockstorage

type Repository interface {
	CreateNewUser(ctx context.Context, login string, passwordHash, salt []byte) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)

	CreatePassword(ctx context.Context, userID string, data []byte) error
	CreateBank(ctx context.Context, userID string, data []byte) error
	CreateTextual(ctx context.Context, userID string, data []byte) error
	CreateReference(ctx context.Context, userID string, id string, filename []byte) error
	CreateData(ctx context.Context, id string, stream pb.Gophkeeper_CreateDataServer) ([]byte, error)

	GetPasswordsByUserID(ctx context.Context, id string) ([][]byte, error)
	GetBankByUserID(ctx context.Context, id string) ([][]byte, error)
	GetTextualByUserID(ctx context.Context, id string) ([][]byte, error)
	GetReferenceByUserID(ctx context.Context, userID string) ([]dto.FileReference, error)
	Download(ctx context.Context, id string, stream pb.Gophkeeper_DownloadServer) error
}

func newRepository(cfg *config.Config, p *pgstorage.PGStorage) Repository {
	return p
}

func Provide() fx.Option {
	return fx.Provide(newRepository)
}
