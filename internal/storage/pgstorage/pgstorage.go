package pgstorage

import (
	"context"

	"github.com/MukizuL/GophKeeper/internal/models"
)

func (s PGStorage) CreateNewUser(ctx context.Context, login, passwordHash string) error {
	_, err := s.conn.Exec(ctx, `INSERT INTO users (login, password_hash) VALUES ($1, $2)`, login, passwordHash)
	if err != nil {
		return err
	}

	return nil
}

func (s PGStorage) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, password_hash FROM users WHERE id = $1`, id).
		Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s PGStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	err := s.conn.QueryRow(ctx, `SELECT id, login, password_hash FROM users WHERE login = $1`, login).
		Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
