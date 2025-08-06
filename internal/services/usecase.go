package services

import (
	"context"
	"errors"

	"github.com/MukizuL/GophKeeper/internal/errs"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s Services) CreateNewUser(ctx context.Context, login, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("failed to hash a password", zap.Error(err))
		return errs.ErrInternalServerError
	}

	err = s.storage.CreateNewUser(ctx, login, string(passwordHash))
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

// Login returns a JWT token and an error
func (s Services) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.storage.GetUserByLogin(ctx, login)
	if err != nil {
		s.logger.Error("failed to get a user by login", zap.String("login", login), zap.Error(err))

		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return "", errs.ErrWrongCredentials
		default:
			return "", errs.ErrInternalServerError
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errs.ErrWrongCredentials
	}

	token, err := s.jwtService.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
