package tests

import (
	ssoV1 "SSO/pkg/proto/sso"
	"SSO/tests/sute"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var appKey = []byte("key")

func TestRegisterLogin(t *testing.T) {
	ctx, st := sute.New(t)

	login := "login"
	pass := "test"

	_, err := st.AuthClient.Register(ctx, &ssoV1.RegisterRequest{
		AppKey:   appKey,
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.Login(ctx, &ssoV1.LoginRequest{
		AppKey:   appKey,
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)

	token := respLogin.GetToken()
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
	assert.InDelta(t, loginTime.Add(st.Cnf.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}
