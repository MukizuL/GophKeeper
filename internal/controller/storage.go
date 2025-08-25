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

func (c Controller) GetPasswords(ctx context.Context, in *pb.GetPasswordsRequest) (*pb.GetPasswordsResponse, error) {
	token, err := helpers.GetToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	data, err := c.services.GetPasswords(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.GetPasswordsResponse{
		Data: data,
	}, nil
}
