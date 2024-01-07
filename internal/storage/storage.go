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
	Delete(ctx context.Context, appId int32, login string) error
	TestOnExist(ctx context.Context, appId int32, login string) (bool, error)
}

type AppsStorage interface {
	Save(ctx context.Context, key []byte) error
	GetByKey(ctx context.Context, key []byte) (models.App, error)
	DeleteByKey(ctx context.Context, key []byte) error
	TestOnExist(ctx context.Context, key []byte) bool
	GetAll(ctx context.Context) ([]*models.App, error)
}

type PermissionsStorage interface {
	Save(ctx context.Context, userId int, value int32) error
	Get(ctx context.Context, userId int) (int32, error)
	Update(ctx context.Context, userId int, value int32) error
	Delete(ctx context.Context, userId int) error
}

type Storage struct {
	UserStorage        UserStorage
	AppStorage         AppsStorage
	PermissionsStorage PermissionsStorage
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
		UserStorage:        mysql.NewUserStorage(db),
		AppStorage:         mysql.NewAppStorage(db),
		PermissionsStorage: mysql.NewPermissionsStorage(db),
	}, nil
}
