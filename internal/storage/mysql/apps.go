package mysql

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage/storageErrors"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AppStorage struct {
	db *sql.DB
}

func NewAppStorage(db *sql.DB) *AppStorage {
	return &AppStorage{
		db: db,
	}
}

func (a *AppStorage) Save(ctx context.Context, key []byte) error {
	if _, err := a.db.ExecContext(ctx, "INSERT INTO apps (secret_key) VALUES (?)", key); err != nil {
		return err
	}
	return nil
}

func (a *AppStorage) GetByKey(ctx context.Context, key []byte) (models.App, error) {
	var app models.App
	if err := a.db.QueryRowContext(ctx, "SELECT * FROM apps WHERE secret_key=?", key).Scan(
		&app.Id, &app.Key,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app, storageErrors.ErrAppNotFound
		}
		return app, err
	}
	return app, nil
}

func (a *AppStorage) DeleteByKey(ctx context.Context, key []byte) error {
	const op = "mysql.AppStorage.DeleteByKey"
	if _, err := a.db.ExecContext(ctx, "DELETE FROM apps WHERE secret_key=?", key); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *AppStorage) TestOnExist(ctx context.Context, key []byte) bool {
	var count int
	_ = a.db.QueryRowContext(ctx, "SELECT COUNT(id) FROM apps WHERE secret_key=?", key).Scan(&count)
	return count != 0
}

func (a *AppStorage) GetAll(ctx context.Context) ([]*models.App, error) {
	const op = "mysql.AppStorage.GetAll"
	var apps []*models.App

	rows, err := a.db.QueryContext(ctx, "SELECT * FROM apps")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	for rows.Next() {
		var app models.App
		if err := rows.Scan(&app.Id, &app.Key); err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}

	return apps, nil
}