package pgstorage

import (
	"context"

	"github.com/MukizuL/GophKeeper/internal/dto"
	"github.com/MukizuL/GophKeeper/internal/models"
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

func (s PGStorage) CreateReference(ctx context.Context, userID string, id string, filename []byte) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO files (id, user_id, filename) VALUES ($1, $2, $3)`, id, userID, filename)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) GetReferenceByUserID(ctx context.Context, userID string) ([]dto.FileReference, error) {
	var out []dto.FileReference
	rows, err := s.conn.Query(ctx, `SELECT id, filename FROM files WHERE user_id = $1 ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var data dto.FileReference
		if err := rows.Scan(&data.ID, &data.Filename); err != nil {
			return nil, err
		}

		out = append(out, data)
	}

	return out, nil
}

func (s PGStorage) GetPasswordsByUserID(ctx context.Context, userID string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM passwords WHERE user_id = $1 ORDER BY id`, userID)
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

func (s PGStorage) GetBankByUserID(ctx context.Context, userID string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM bank WHERE user_id = $1 ORDER BY id`, userID)
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

func (s PGStorage) GetTextualByUserID(ctx context.Context, userID string) ([][]byte, error) {
	var out [][]byte
	rows, err := s.conn.Query(ctx, `SELECT data FROM textual WHERE user_id = $1 ORDER BY id`, userID)
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
