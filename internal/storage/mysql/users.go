package mysql

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage/storageErrors"
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	const op = "userStorage.Save"
	if _, err := u.db.ExecContext(ctx, "INSERT INTO users (login, password, app_id) VALUES (?, ?, ?)", login, passwordHash, appid); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserStorage) Get(ctx context.Context, appId int32, login string) (models.User, error) {
	const op = "userStorage.Get"
	var user models.User

	if err := u.db.QueryRowContext(ctx, "SELECT * FROM users WHERE app_id=? AND login=?", appId, login).Scan(
		&user.Id, &user.AppId, &user.Login, &user.PasswordHash,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, storageErrors.ErrUserNotFound
		}
		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (u *UserStorage) Delete(ctx context.Context, appId int32, login string) error {
	const op = "userStorage.Delete"
	if _, err := u.db.ExecContext(ctx, "DELETE FROM users WHERE app_id=? AND login=?", appId, login); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserStorage) UpdateLogin(ctx context.Context, appId int32, login string, newLogin string) error {
	const op = "userStorage.UpdateLogin"
	if _, err := u.db.ExecContext(ctx, "UPDATE users SET login=? WHERE app_id=? AND login=?;", newLogin, appId, login); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (u *UserStorage) TestOnExist(ctx context.Context, appId int32, login string) (bool, error) {
	const op = "userStorage.TestOnExist"
	var count int
	if err := u.db.QueryRowContext(ctx, "SELECT COUNT(id) FROM users WHERE app_id=? AND login=?", appId, login).Scan(&count); err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	return count != 0, nil
}
