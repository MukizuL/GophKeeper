package controller

import (
	"context"
	"errors"
	"unicode/utf8"

	"github.com/MukizuL/GophKeeper/internal/errs"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c Controller) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if utf8.RuneCountInString(in.Login) < 3 {
		return nil, status.Error(codes.InvalidArgument, "login must be at least 3 characters")
	}
	if utf8.RuneCountInString(in.Login) > 255 {
		return nil, status.Error(codes.InvalidArgument, "login must be at most 255 characters")
	}
	if utf8.RuneCountInString(in.Password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters")
	}
	if utf8.RuneCountInString(in.Password) > 36 {
		return nil, status.Error(codes.InvalidArgument, "password must be at most 36 characters")
	}
	if len(in.Password) > 72 {
		return nil, status.Errorf(codes.InvalidArgument, "password is %d characters but is longer than 72 bytes", utf8.RuneCountInString(in.Password))
	}

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
