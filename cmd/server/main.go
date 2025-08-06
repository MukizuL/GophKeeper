package main

import (
	"fmt"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/controller"
	"github.com/MukizuL/GophKeeper/internal/interceptor"
	jwtService "github.com/MukizuL/GophKeeper/internal/jwt"
	"github.com/MukizuL/GophKeeper/internal/migration"
	"github.com/MukizuL/GophKeeper/internal/server"
	"github.com/MukizuL/GophKeeper/internal/services"
	"github.com/MukizuL/GophKeeper/internal/storage"
	"github.com/MukizuL/GophKeeper/internal/storage/pgstorage"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n\n", buildCommit)

	fx.New(
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
		createApp(),
		fx.Invoke(func(*grpc.Server) {}),
	).Run()
}

func createApp() fx.Option {
	return fx.Options(
		config.Provide(),
		fx.Provide(zap.NewDevelopment),
		server.Provide(),
		jwtService.Provide(),
		interceptor.Provide(),

		controller.Provide(),
		services.Provide(),

		pgstorage.Provide(),
		storage.Provide(),
		migration.Provide(),
	)
}
