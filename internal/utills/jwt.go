package utills

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sso_service_grps/internal/domain/models"
	"time"
)

func GenerateToken(user models.User, app models.App, duration time.Duration) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  user.Email,
		"userId": user.ID,
		"exp":    time.Now().Add(duration).Unix(),
		"appId":  app.ID,
	})
	return token.SignedString([]byte(app.Secret))
}

func VerifyToken(token string, app models.App) (int, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("unexpected signed method")
		}
		return []byte(app.Secret), nil
	})
	if err != nil {
		return 0, errors.New("could not parse token")
	}
	if !parsedToken.Valid {
		return 0, errors.New("invalid token")
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userId := claims["userId"].(float64)

	return int(userId), nil
}
