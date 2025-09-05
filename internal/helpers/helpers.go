package helpers

import (
	"context"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/MukizuL/GophKeeper/internal/ctxutil"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
	"google.golang.org/grpc/metadata"
)

const chunkSize = 32_796

// GetToken extracts access token from incoming context
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

func GetUserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(ctxutil.UserIDContextKey).(string)
	if !ok {
		return "", errors.New("userID is not a string")
	}
	return userID, nil
}

func PrepareFile(fullPath string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return nil, err
	}
	return os.Create(fullPath)
}

func ReceiveChunks(stream pb.Gophkeeper_CreateDataServer, f *os.File) ([]byte, error) {
	var filepath []byte
	first := true

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if first {
			filepath = chunk.Filename
			first = false
		}

		if _, err := f.Write(chunk.Chunk); err != nil {
			return nil, err
		}
	}
	return filepath, nil
}

func SendChunks(stream pb.Gophkeeper_DownloadServer, f *os.File) error {
	buf := make([]byte, chunkSize) // 32MB + 28 bytes

	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return errors.New("could not read file")
		}

		req := &pb.DownloadResponse{
			Chunk: buf[:n],
		}

		if err := stream.Send(req); err != nil {
			return errors.New("could not send data")
		}
	}
}
