package server

import (
	"context"
	"fmt"
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
	"google.golang.org/grpc/credentials"
)

type GRPCFxIn struct {
	fx.In

	Lc          fx.Lifecycle
	Ctrl        *controller.Controller
	Cfg         *config.Config
	Logger      *zap.Logger
	Interceptor *interceptor.Service
	Storage     storage.Repository
	Migrator    *migration.Migrator
}

func newGRPCServer(in GRPCFxIn) (*grpc.Server, error) {
	ln, err := net.Listen("tcp", in.Cfg.GRPCPort)
	if err != nil {
		return nil, err
	}

	var s *grpc.Server
	if !in.Cfg.TLS {
		in.Logger.Info("Starting GRPC server w/out TLS", zap.String("port", in.Cfg.GRPCPort))
		s = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				in.Interceptor.Logger,
				in.Interceptor.Auth,
			),
		)
	} else {
		in.Logger.Info("Starting GRPC server with TLS", zap.String("port", in.Cfg.GRPCPort))
		creds, err := credentials.NewServerTLSFromFile(in.Cfg.Cert, in.Cfg.PK)
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}

		s = grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				in.Interceptor.Logger,
				in.Interceptor.Auth,
			),
			grpc.Creds(creds),
		)
	}

	pb.RegisterGophkeeperServer(s, in.Ctrl)

	in.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
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
