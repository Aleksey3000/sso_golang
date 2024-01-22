package jwt

import (
	"SSO/internal/domain/models"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrExpired = errors.New("token has expired")
)

func NewToken(user models.User, app models.App, TTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(TTL).Unix()

	tokenStr, err := token.SignedString(app.Key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ParseToken(strToken string, key []byte) (string, error) {
	token, err := jwt.Parse(strToken, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	login := claims["login"].(string)
	exp := int64(claims["exp"].(float64))
	if time.Now().Unix() > exp {
		return "", ErrExpired
	}
	return login, err
}
