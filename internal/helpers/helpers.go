package helpers

import (
	"context"
	"crypto/rand"
	"errors"
	"io"

	"google.golang.org/grpc/metadata"
)

func GetToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("metadata is not provided")
	}

	vals := md.Get("access-token")

	if len(vals) == 0 || vals[0] == "" {
		return "", errors.New("access token is missing")
	}

	return vals[0], nil
}

// GenerateSalt creates a random salt for key derivation.
func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}
