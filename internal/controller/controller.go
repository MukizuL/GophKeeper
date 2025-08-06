package controller

import (
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/MukizuL/GophKeeper/internal/services"
	"github.com/MukizuL/GophKeeper/internal/storage"
	"go.uber.org/fx"
)

type Controller struct {
	storage  storage.Repo
	services *services.Services
	pb.UnimplementedGophkeeperServer
}

func newController(storage storage.Repo, services *services.Services) *Controller {
	return &Controller{
		storage:  storage,
		services: services,
	}
}

func Provide() fx.Option {
	return fx.Provide(newController)
}
