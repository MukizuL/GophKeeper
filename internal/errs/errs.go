package errs

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrDuplicateLogin   = errors.New("this login is already used")
	ErrUserNotFound     = errors.New("user with this login is not found")
	ErrWrongCredentials = errors.New("wrong credentials")

	ErrNotAuthorized           = errors.New("invalid token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrUserMismatch            = errors.New("user tried to delete not owned urls")
	ErrGone                    = errors.New("url was marked as deleted")
	ErrSigningToken            = errors.New("error signing token")
	ErrRefreshingToken         = errors.New("error refreshing token")
	ErrNoCert                  = errors.New("no certificate provided")
	ErrNoPK                    = errors.New("no private key provided")
	ErrInternalServerError     = errors.New("internal server error")
)

func TransformPGErrors(err error) error {
	var pgErr *pgconn.PgError

	switch {
	case errors.As(err, &pgErr):
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrDuplicateLogin
		}
	case errors.Is(err, pgx.ErrNoRows):
		return ErrUserNotFound
	default:
		return ErrInternalServerError
	}

	return ErrInternalServerError
}
