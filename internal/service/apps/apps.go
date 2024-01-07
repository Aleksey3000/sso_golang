package apps

import (
	"SSO/internal/domain/models"
	"SSO/internal/storage"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"strconv"
	"sync"
	"time"
)

type Apps struct {
	l           *slog.Logger
	appsStorage storage.AppsStorage
}

func New(l *slog.Logger, appsStorage storage.AppsStorage) *Apps {
	return &Apps{
		l:           l,
		appsStorage: appsStorage,
	}
}

func (a *Apps) NewApp(ctx context.Context) ([]byte, error) {
	key := GenerateUniqueString()
	if err := a.appsStorage.Save(ctx, key); err != nil {
		a.l.Error(err.Error())
		return nil, err
	}
	return key, nil
}

func (a *Apps) DeleteApp(ctx context.Context, key []byte) error {
	return a.appsStorage.DeleteByKey(ctx, key)
}

func (a *Apps) TestOnExist(ctx context.Context, key []byte) bool {
	return a.appsStorage.TestOnExist(ctx, key)
}

func (a *Apps) GetAll(ctx context.Context) ([]*models.App, error) {
	return a.appsStorage.GetAll(ctx)
}

var mu sync.Mutex

func GenerateUniqueString() []byte {
	mu.Lock()
	num := time.Now().UnixNano()
	mu.Unlock()
	hash := sha256.New()

	hash.Write([]byte(strconv.Itoa(int(num))))
	b := hash.Sum(nil)

	buf := make([]byte, len(b)*2)
	hex.Encode(buf, b)
	return buf
}
