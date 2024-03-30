package util

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/splorg/go-auth/internal/model"
)

func GenerateJWT(user *model.User, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add((time.Hour * 24) * 7).Unix(),
    IssuedAt: time.Now().Unix(),
    Subject: user.ID,
	})

	return claims.SignedString([]byte(secret))
}

func GenerateCookie(token string) *fiber.Cookie {
	return &fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add((time.Hour * 24) * 7),
		HTTPOnly: true,
	}
}
