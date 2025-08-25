package controller

import (
	"context"
	"errors"

	"github.com/MukizuL/GophKeeper/internal/errs"
	"github.com/MukizuL/GophKeeper/internal/helpers"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c Controller) CreatePassword(ctx context.Context, in *pb.CreatePasswordRequest) (*pb.CreatePasswordResponse, error) {
	token, err := helpers.GetToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = c.services.CreatePassword(ctx, token, in.Data)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.CreatePasswordResponse{}, nil
}
