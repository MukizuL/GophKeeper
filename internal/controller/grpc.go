package controller

import (
	"context"
	"errors"

	"github.com/MukizuL/GophKeeper/internal/errs"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c Controller) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := c.services.CreateNewUser(ctx, in.Login, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateLogin):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &pb.RegisterResponse{}, nil
}

func (c Controller) Authorize(ctx context.Context, in *pb.AuthRequest) (*pb.AuthResponse, error) {
	token, err := c.services.Login(ctx, in.Login, in.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, errs.ErrWrongCredentials):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		case errors.Is(err, errs.ErrSigningToken):
			return nil, status.Error(codes.Internal, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	header := metadata.Pairs("access-token", token)
	err = grpc.SetHeader(ctx, header)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AuthResponse{}, nil
}
