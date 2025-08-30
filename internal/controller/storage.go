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
	data, err := c.services.GetPasswords(ctx)
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

func (c Controller) GetBank(ctx context.Context, in *pb.GetBankRequest) (*pb.GetBankResponse, error) {
	data, err := c.services.GetBank(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.GetBankResponse{
		Data: data,
	}, nil
}

func (c Controller) GetText(ctx context.Context, in *pb.GetTextRequest) (*pb.GetTextResponse, error) {
	data, err := c.services.GetText(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.GetTextResponse{
		Data: data,
	}, nil
}

func (c Controller) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
	data, err := c.services.GetData(ctx)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, "Internal Server Error")
		}
	}

	return &pb.GetDataResponse{
		File: data,
	}, nil
}
func (c Controller) Download(in *pb.DownloadRequest, stream pb.Gophkeeper_DownloadServer) error {
	ctx := stream.Context()

	token, err := helpers.GetToken(ctx)
	if err != nil {
		return status.Error(codes.Unauthenticated, err.Error())
	}

	err = c.services.Download(stream.Context(), token, in.Id, stream)
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
