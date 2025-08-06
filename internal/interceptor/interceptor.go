package interceptor

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/ctxutil"
	"github.com/MukizuL/GophKeeper/internal/errs"
	jwtService "github.com/MukizuL/GophKeeper/internal/jwt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Service struct {
	jwtService jwtService.ServiceI
	cfg        *config.Config
	logger     *zap.Logger
}

func newService(jwtService jwtService.ServiceI, cfg *config.Config, logger *zap.Logger) *Service {
	return &Service{
		jwtService: jwtService,
		cfg:        cfg,
		logger:     logger,
	}
}

func Provide() fx.Option {
	return fx.Provide(newService)
}

func (s Service) Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	resp, err := handler(ctx, req)

	duration := time.Since(start)

	var client string
	pr, ok := peer.FromContext(ctx)
	if ok {
		client = pr.Addr.String()
	}

	s.logger.Info("GRPC request",
		zap.String("method", info.FullMethod),
		zap.String("client", client),
		zap.Duration("time", duration))

	return resp, err
}

func (s Service) Auth(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	routes := []string{
		"/shortener.Shortener/CreateGRPC",
		"/shortener.Shortener/CreateBatchGRPC",
		"/shortener.Shortener/GetUserURLsGRPC",
		"/shortener.Shortener/DeleteGRPC",
	}

	if !slices.Contains(routes, info.FullMethod) {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	var userID string
	var err error

	tokens := md.Get("access-token")

	if len(tokens) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "no access token")
	} else {
		userID, err = s.jwtService.ValidateToken(tokens[0])
	}
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotAuthorized), errors.Is(err, errs.ErrUnexpectedSigningMethod):
			return nil, status.Errorf(codes.Unauthenticated, "%s", err.Error())
		case errors.Is(err, errs.ErrSigningToken):
			return nil, status.Errorf(codes.Internal, "%s", err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "%s", err.Error())
		}
	}

	newCtx := context.WithValue(ctx, ctxutil.UserIDContextKey, userID)

	return handler(newCtx, req)
}

func (s Service) Recovery(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Error("panic recovered",
				zap.Any("panic", r),
				zap.String("method", info.FullMethod),
				zap.Stack("stack"),
			)
			err = status.Error(codes.Internal, "internal server error")
		}
	}()

	return handler(ctx, req)
}
