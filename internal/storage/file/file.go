package file

import (
	"context"
	"os"
	"path/filepath"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/helpers"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"go.uber.org/fx"
)

type Storage struct {
	cfg *config.Config
}

func newStorage(cfg *config.Config) *Storage {
	return &Storage{
		cfg: cfg,
	}
}

func Provide() fx.Option {
	return fx.Provide(newStorage)
}

func (s Storage) CreateData(ctx context.Context, id string, stream pb.Gophkeeper_CreateDataServer) ([]byte, error) {
	fullPath := filepath.Join(s.cfg.Filepath, id)

	f, err := helpers.PrepareFile(fullPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fPath, err := helpers.ReceiveChunks(stream, f)
	if err != nil {
		return nil, err
	}

	return fPath, stream.SendAndClose(&pb.CreateDataResponse{})
}

func (s Storage) Download(ctx context.Context, id string, stream pb.Gophkeeper_DownloadServer) error {
	fullPath := filepath.Join(s.cfg.Filepath, id)

	f, err := os.Open(fullPath)
	if err != nil {
		return err
	}

	return helpers.SendChunks(stream, f)
}
