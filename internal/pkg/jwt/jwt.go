package jwt

import (
	"SSO/internal/domain/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.App, TTL time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.Id
	claims["login"] = user.Login
	claims["exp"] = time.Now().Add(TTL).Unix()
	claims["app_id"] = app.Id

	tokenStr, err := token.SignedString(app.Key)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}
