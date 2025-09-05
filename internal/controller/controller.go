package controller

import (
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/MukizuL/GophKeeper/internal/services"
	"go.uber.org/fx"
)

type Controller struct {
	services services.ServicesI
	pb.UnimplementedGophkeeperServer
}

func newController(services services.ServicesI) *Controller {
	return &Controller{
		services: services,
	}
}

func Provide() fx.Option {
	return fx.Provide(newController)
}
