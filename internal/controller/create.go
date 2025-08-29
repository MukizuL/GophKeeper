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

func (c Controller) CreateBank(ctx context.Context, in *pb.CreateBankRequest) (*pb.CreateBankResponse, error) {
	token, err := helpers.GetToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = c.services.CreateBank(ctx, token, in.Data)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.CreateBankResponse{}, nil
}

func (c Controller) CreateTextual(ctx context.Context, in *pb.CreateTextRequest) (*pb.CreateTextResponse, error) {
	token, err := helpers.GetToken(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	err = c.services.CreateTextual(ctx, token, in.Data)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.CreateTextResponse{}, nil
}

func (c Controller) CreateData(stream pb.Gophkeeper_CreateDataServer) error {
	ctx := stream.Context()

	token, err := helpers.GetToken(ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	err = c.services.CreateData(ctx, token, stream)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return status.Error(codes.Unauthenticated, err.Error())
		default:
			return status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return nil
}
