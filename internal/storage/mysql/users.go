package mysql

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage/storageErrors"
	"context"
	"database/sql"
	"errors"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (u *UserStorage) Save(ctx context.Context, appid int32, login string, passwordHash []byte) error {
	if _, err := u.db.ExecContext(ctx, "INSERT INTO users (login, password, app_id) VALUES (?, ?, ?)", login, passwordHash, appid); err != nil {
		return err
	}
	return nil
}

func (u *UserStorage) Get(ctx context.Context, appId int32, login string) (models.User, error) {
	var user models.User

	if err := u.db.QueryRowContext(ctx, "SELECT * FROM users WHERE app_id=? AND login=?", appId, login).Scan(
		&user.Id, &user.AppId, &user.Login, &user.PasswordHash,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, storageErrors.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}
