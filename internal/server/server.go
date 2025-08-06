package server

import (
	"context"
	"net"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/controller"
	"github.com/MukizuL/GophKeeper/internal/interceptor"
	"github.com/MukizuL/GophKeeper/internal/migration"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/MukizuL/GophKeeper/internal/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type GRPCFxIn struct {
	fx.In

	Lc          fx.Lifecycle
	Ctrl        *controller.Controller
	Cfg         *config.Config
	Logger      *zap.Logger
	Interceptor *interceptor.Service
	Storage     storage.Repo
	Migrator    *migration.Migrator
}

func newGRPCServer(in GRPCFxIn) (*grpc.Server, error) {
	ln, err := net.Listen("tcp", in.Cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			in.Interceptor.Logger,
			in.Interceptor.Auth,
		),
	)

	pb.RegisterGophkeeperServer(s, in.Ctrl)

	in.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			in.Logger.Info("Starting GRPC server", zap.String("port", in.Cfg.GRPCPort))
			go s.Serve(ln)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			go s.GracefulStop()

			return nil
		},
	})

	return s, nil
}

func Provide() fx.Option {
	return fx.Provide(newGRPCServer)
}
