package storage

import (
	"SSO/internal/config"
	"SSO/internal/domain/models"
	"SSO/internal/storage/mysql"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type UserStorage interface {
	Save(ctx context.Context, appId int32, login string, passwordHash []byte) error
	Get(ctx context.Context, appId int32, login string) (models.User, error)
}

type AppStorage interface {
	Save(ctx context.Context, key []byte) error
	GetByKey(ctx context.Context, key []byte) (models.App, error)
}

type Storage struct {
	UserStorage UserStorage
	AppStorage  AppStorage
}

func New(cnf *config.DBConfig) (*Storage, error) {
	const op = "storage.New"
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s", cnf.User, cnf.Password, cnf.Server, cnf.DBName))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{
		UserStorage: mysql.NewUserStorage(db),
		AppStorage:  mysql.NewAppStorage(db),
	}, nil
}
