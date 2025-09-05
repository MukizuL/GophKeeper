package services

import (
	"context"
	"errors"

	"github.com/MukizuL/GophKeeper/internal/ctxutil"
	"github.com/MukizuL/GophKeeper/internal/errs"
	"github.com/MukizuL/GophKeeper/internal/helpers"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s Services) CreateNewUser(ctx context.Context, login, password string) error {
	salt, err := helpers.GenerateSalt(16)
	if err != nil {
		s.logger.Error("failed to generate salt", zap.Error(err))
		return errs.ErrInternalServerError
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash a password", zap.Error(err))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreateNewUser(ctx, login, passwordHash, salt)
	if err != nil {
		s.logger.Error("failed to create a new user", zap.Error(err))

		var pgErr *pgconn.PgError

		switch {
		case errors.As(err, &pgErr):
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				return errs.ErrDuplicateLogin
			default:
				return errs.ErrInternalServerError
			}
		default:
			return errs.ErrInternalServerError
		}
	}

	return nil
}

// Login returns a JWT token, salt and an error
func (s Services) Login(ctx context.Context, login, password string) (string, []byte, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("failed to get a user by login", zap.String("login", login), zap.Error(err))

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return "", nil, errs.ErrWrongCredentials
		default:
			return "", nil, errs.ErrInternalServerError
		}
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return "", nil, errs.ErrWrongCredentials
	}

	token, err := s.jwtService.CreateToken(user.ID)
	if err != nil {
		return "", nil, err
	}

	return token, user.Salt, nil
}

func (s Services) CreatePassword(ctx context.Context, data []byte) error {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreatePassword(ctx, userID, data)
	if err != nil {
		s.logger.Error("failed to create a new password", zap.String("userID", userID), zap.Error(err))
		return errs.ErrInternalServerError
	}

	return nil
}

func (s Services) CreateBank(ctx context.Context, data []byte) error {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreateBank(ctx, userID, data)
	if err != nil {
		s.logger.Error("failed to create a new bank card", zap.String("userID", userID), zap.Error(err))
		return errs.ErrInternalServerError
	}

	return nil
}

func (s Services) CreateTextual(ctx context.Context, data []byte) error {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreateTextual(ctx, userID, data)
	if err != nil {
		s.logger.Error("failed to create a new text", zap.String("userID", userID), zap.Error(err))
		return errs.ErrInternalServerError
	}

	return nil
}

func (s Services) CreateData(ctx context.Context, token string, stream pb.Gophkeeper_CreateDataServer) error {
	userID, err := s.jwtService.ValidateToken(token)
	if err != nil {
		s.logger.Error("failed to validate token", zap.String("token", token), zap.Error(err))
		return err
	}

	id := uuid.New()

	filename, err := s.storage.CreateData(ctx, id.String(), stream)
	if err != nil {
		s.logger.Error("failed to create a new data",
			zap.String("id", id.String()),
			zap.String("userID", userID),
			zap.Error(err))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreateReference(ctx, userID, id.String(), filename)
	if err != nil {
		s.logger.Error("failed to create a new data reference",
			zap.String("id", id.String()),
			zap.String("userID", userID),
			zap.Error(err))
		return errs.ErrInternalServerError
	}

	return nil
}

func (s Services) GetPasswords(ctx context.Context) ([][]byte, error) {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return nil, errs.ErrInternalServerError
	}

	response, err := s.storage.GetPasswordsByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get a user's passwords", zap.String("userID", userID), zap.Error(err))
		return nil, errs.ErrInternalServerError
	}

	return response, nil
}

func (s Services) GetBank(ctx context.Context) ([][]byte, error) {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return nil, errs.ErrInternalServerError
	}

	response, err := s.storage.GetBankByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get a user's bank cards", zap.String("userID", userID), zap.Error(err))
		return nil, errs.ErrInternalServerError
	}

	return response, nil
}

func (s Services) GetText(ctx context.Context) ([][]byte, error) {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return nil, errs.ErrInternalServerError
	}

	response, err := s.storage.GetTextualByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get a user's texts", zap.String("userID", userID), zap.Error(err))
		return nil, errs.ErrInternalServerError
	}

	return response, nil
}

func (s Services) GetData(ctx context.Context) ([]*pb.File, error) {
	userID, err := helpers.GetUserIDFromContext(ctx)
	if err != nil {
		s.logger.Error(err.Error(), zap.Any("userID", ctx.Value(ctxutil.UserIDContextKey)))
		return nil, errs.ErrInternalServerError
	}

	refs, err := s.storage.GetReferenceByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to get a user's texts", zap.String("userID", userID), zap.Error(err))
		return nil, errs.ErrInternalServerError
	}

	var response []*pb.File
	for _, v := range refs {
		response = append(response, &pb.File{
			Id:       v.ID,
			Filename: v.Filename,
		})
	}

	return response, nil
}

func (s Services) Download(ctx context.Context, token, id string, stream pb.Gophkeeper_DownloadServer) error {
	_, err := s.jwtService.ValidateToken(token)
	if err != nil {
		s.logger.Error("failed to validate token", zap.String("token", token), zap.Error(err))
		return err
	}

	err = s.storage.Download(ctx, id, stream)
	if err != nil {
		s.logger.Error("failed to download a file", zap.String("id", id), zap.Error(err))
		return errs.ErrInternalServerError
	}

	return nil
}
