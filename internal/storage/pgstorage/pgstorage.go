package pgstorage

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/MukizuL/GophKeeper/internal/models"
	pb "github.com/MukizuL/GophKeeper/internal/proto"
)

func (s PGStorage) CreateNewUser(ctx context.Context, login string, passwordHash, salt []byte) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO users (login, password_hash, kdf_salt) VALUES ($1, $2, $3)`, login, passwordHash, salt)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, password_hash, kdf_salt FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Login, &user.Password, &user.Salt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s PGStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, password_hash, kdf_salt FROM users WHERE login = $1`, login).
		Scan(&user.ID, &user.Login, &user.Password, &user.Salt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s PGStorage) CreatePassword(ctx context.Context, userID string, data []byte) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO passwords (user_id, data) VALUES ($1, $2)`, userID, data)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) CreateBank(ctx context.Context, userID string, data []byte) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO bank (user_id, data) VALUES ($1, $2)`, userID, data)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) CreateTextual(ctx context.Context, userID string, data []byte) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO textual (user_id, data) VALUES ($1, $2)`, userID, data)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) CreateReference(ctx context.Context, userID string, id, filename string) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO files (id, user_id, filename) VALUES ($1, $2, $3)`, id, userID, filename)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) CreateData(ctx context.Context, id string, stream pb.Gophkeeper_CreateDataServer) (string, error) {
	fullPath := filepath.Join(s.cfg.Filepath, id)

	err := os.MkdirAll(filepath.Dir(fullPath), 0755)
	if err != nil {
		return "", err
	}

	f, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	var filepath string
	first := true

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		if first {
			filepath = chunk.Filename
			first = false
		}

		if _, err := f.Write(chunk.Chunk); err != nil {
			return "", err
		}
	}

	return filepath, stream.SendAndClose(&pb.CreateDataResponse{})
}

func (s PGStorage) GetPasswordsByUserID(ctx context.Context, id string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM passwords WHERE user_id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		out = append(out, data)
	}

	return out, nil
}

func (s PGStorage) GetBankByUserID(ctx context.Context, id string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM bank WHERE user_id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		out = append(out, data)
	}

	return out, nil
}

func (s PGStorage) GetTextualByUserID(ctx context.Context, id string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM textual WHERE user_id = $1 ORDER BY id`, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			return nil, err
		}

		out = append(out, data)
	}

	return out, nil
}
