package main

import (
	"context"
	"fmt"
	"time"

	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Register(login, password string) error {
	err := validateLogin(login)
	if err != nil {
		return err
	}

	err = validatePassword(password)
	if err != nil {
		return err
	}

	req := pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = conn.Register(ctx, &req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.FailedPrecondition:
				return fmt.Errorf("user with same login already exists: %s", e.Message())
			case codes.Internal:
				return fmt.Errorf("server error: %s", e.Message())
			default:
				return fmt.Errorf("unknown error: %s", e.Message())
			}
		}
	}

	return nil
}

func Login(login, password string) (string, error) {
	err := validateLogin(login)
	if err != nil {
		return "", err
	}

	err = validatePassword(password)
	if err != nil {
		return "", err
	}

	req := pb.AuthRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var header metadata.MD
	_, err = conn.Authorize(ctx, &req, grpc.Header(&header))
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return "", fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.Unauthenticated:
				return "", fmt.Errorf("%s", e.Message())
			case codes.Internal:
				return "", fmt.Errorf("server error: %s", e.Message())
			default:
				return "", fmt.Errorf("unknown error: %s", e.Message())
			}
		}
	}

	tokens := header.Get("access-token")
	if len(tokens) == 0 {
		return "", fmt.Errorf("access token is missing")
	}

	return tokens[0], nil
}
