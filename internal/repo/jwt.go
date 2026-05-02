package repo

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/logservice/internal/model"
)

func GenerateAccessToken(user model.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	claims := model.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(), // sub
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(user model.User) (string, error) {
	secret := os.Getenv("JWT_REFRESH_SECRET")

	claims := model.Claims{
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(), // sub
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}