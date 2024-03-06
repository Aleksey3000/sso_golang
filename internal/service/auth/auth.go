package auth

import (
	"SSO/internal/domain/models"
	"SSO/internal/pkg/jwt"
	"SSO/internal/storage"
	"SSO/internal/storage/storageErrors"
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserIsExist        = errors.New("user is exist")
)

type AppsProvider interface {
	GetByKey(ctx context.Context, key []byte) (models.App, error)
}

type PermDeleter interface {
	Delete(ctx context.Context, userId int64) error
}

type Auth struct {
	l            *slog.Logger
	userStorage  storage.UserStorage
	appsProvider AppsProvider
	perm         PermDeleter
	tokenTTL     time.Duration
}

func New(l *slog.Logger, userStorage storage.UserStorage, appProvider AppsProvider, perm PermDeleter, tokenTTL time.Duration) *Auth {
	return &Auth{
		l:            l,
		userStorage:  userStorage,
		appsProvider: appProvider,
		tokenTTL:     tokenTTL,
		perm:         perm,
	}
}

func (a *Auth) Register(ctx context.Context, appKey []byte, login string, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		return err
	}
	if _, err := a.userStorage.Get(ctx, app.Id, login); err == nil {
		a.l.Error("error register", Err(ErrUserIsExist))
		return ErrUserIsExist
	}
	if err := a.userStorage.Save(ctx, app.Id, login, passHash); err != nil {
		return err
	}
	a.l.Info("register user %s", login)
	return nil
}

func (a *Auth) Login(ctx context.Context, appKey []byte, login string, password string) (string, error) {
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		return "", err
	}

	user, err := a.userStorage.Get(ctx, app.Id, login)
	if err != nil {
		if errors.Is(err, storageErrors.ErrUserNotFound) {
			a.l.Info("user %s, app:%d not found", login, app.Id)
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.l.Error("failed generate token", Err(err))
		return "", err
	}
	return token, nil
}

func (a *Auth) DeleteUser(ctx context.Context, appKey []byte, login string) error {
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		a.l.Error("failed get app", Err(err))
		return err
	}
	user, err := a.userStorage.Get(ctx, app.Id, login)
	if err != nil {
		a.l.Error("failed get user", Err(err))
		return err
	}
	if err := a.userStorage.Delete(ctx, app.Id, login); err != nil {
		a.l.Error("failed delete user", Err(err))
		return err
	}
	if err := a.perm.Delete(ctx, user.Id); err != nil {
		a.l.Error("failed delete permission", Err(err))
		return err
	}
	return nil
}

func (a *Auth) TestOnExist(ctx context.Context, appKey []byte, login string) bool {
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		a.l.Error("failed get app", Err(err))
		return false
	}
	exist, err := a.userStorage.TestOnExist(ctx, app.Id, login)
	if err != nil {
		a.l.Error("failed test user on exist", Err(err))
		return false
	}
	return exist
}

func (a *Auth) UpdateLogin(ctx context.Context, appKey []byte, login string, newLogin string) error {
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		a.l.Error("failed get app", Err(err))
		return err
	}
	if _, err := a.userStorage.Get(ctx, app.Id, newLogin); err == nil {
		a.l.Error("error update", Err(ErrUserIsExist))
		return ErrUserIsExist
	}
	if err := a.userStorage.UpdateLogin(ctx, app.Id, login, newLogin); err != nil {
		a.l.Error("failed update login", Err(err))
		return err
	}
	return nil
}

func (a *Auth) GetUserId(ctx context.Context, appKey []byte, login string) (int64, error) {
	app, err := a.appsProvider.GetByKey(ctx, appKey)
	if err != nil {
		a.l.Error("failed get app", Err(err))
		return 0, err
	}
	user, err := a.userStorage.Get(ctx, app.Id, login)
	if err != nil {
		a.l.Error("failed test user on exist", Err(err))
		return 0, err
	}
	return user.Id, nil
}

func (a *Auth) ParseToken(ctx context.Context, appKey []byte, token string) (string, error) {
	login, err := jwt.ParseToken(token, appKey)
	if err != nil {
		a.l.Warn(err.Error())
		return "", err
	}
	return login, nil
}

func Err(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}
