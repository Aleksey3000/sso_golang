package tests

import (
	"SSO/pkg/AuthClient"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	client, err := AuthClient.New("127.0.0.1", "8080", "key")
	require.NoError(t, err)
	ctx := context.Background()

	const (
		login = "test"
		pass  = "pass"
	)
	var (
		tokenTTL = time.Hour * 1
	)

	err = client.Register(ctx, login, pass)
	require.NoError(t, err)

	token, err := client.Login(ctx, login, pass)

	require.NotEmpty(t, token) // Проверяем, что он не пустой

	// Отмечаем время, в которое бы выполнен логин.
	// Это понадобится для проверки TTL токена
	loginTime := time.Now()

	// Парсим и валидируем токен
	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return appKey, nil
	})
	// Если ключ окажется невалидным, мы получим соответствующую ошибку
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, login, claims["login"].(string))

	const deltaSeconds = 1

	// Проверяем, что TTL токена примерно соответствует нашим ожиданиям.
	assert.InDelta(t, loginTime.Add(tokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)

	// delete user

	err = client.DeleteUser(ctx, login)
	require.NoError(t, err)

}
