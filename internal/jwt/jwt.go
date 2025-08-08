package jwt

import (
	"fmt"
	"time"

	"github.com/MukizuL/GophKeeper/internal/config"
	"github.com/MukizuL/GophKeeper/internal/errs"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/fx"
)

//go:generate mockgen -source=jwt.go -destination=mocks/jwt.go -package=mockjwt

type ServiceI interface {
	ValidateToken(token string) (string, error)
	CreateToken(userID string) (string, error)
}

type Service struct {
	key []byte
}

func newJWTService(cfg *config.Config) ServiceI {
	return &Service{
		key: []byte(cfg.MasterPassword),
	}
}

func Provide() fx.Option {
	return fx.Provide(newJWTService)
}

// ValidateToken returns parsed userID and an error
func (s *Service) ValidateToken(token string) (string, error) {
	var claims jwt.RegisteredClaims
	accessToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", errs.ErrUnexpectedSigningMethod, token.Header["alg"])
		}

		return s.key, nil
	})
	if err != nil {
		return "", errs.ErrNotAuthorized
	}

	if !accessToken.Valid {
		return "", errs.ErrNotAuthorized
	}

	return claims.Subject, nil
}

// CreateToken returns a new token and an error
func (s *Service) CreateToken(userID string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(3600 * time.Second)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	accessTokenSigned, err := accessToken.SignedString(s.key)
	if err != nil {
		return "", errs.ErrSigningToken
	}

	return accessTokenSigned, nil
}
