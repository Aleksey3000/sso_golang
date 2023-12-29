package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type PermissionsStorage struct {
	db *sql.DB
}

func NewPermissionsStorage(db *sql.DB) *PermissionsStorage {
	return &PermissionsStorage{
		db: db,
	}
}

func (p *PermissionsStorage) Save(ctx context.Context, userId int, value int32) error {
	const op = "PermissionsStorage.Save"
	if _, err := p.db.ExecContext(ctx, "INSERT INTO permissions (user_id, permission) VALUES (?, ?)", userId, value); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *PermissionsStorage) Get(ctx context.Context, userId int) (int32, error) {
	const op = "PermissionsStorage.Get"
	var perm int32
	if err := p.db.QueryRowContext(ctx, "SELECT permission FROM permissions WHERE user_id=?", userId).Scan(&perm); err != nil {
		return perm, fmt.Errorf("%s: %w", op, err)
	}
	return perm, nil
}

func (p *PermissionsStorage) Update(ctx context.Context, userId int, value int32) error {
	const op = "PermissionsStorage.Update"
	//todo
	return nil
}

func (p *PermissionsStorage) Delete(ctx context.Context, userId int) error {
	const op = "PermissionsStorage.Delete"
	// todo
	return nil
}
