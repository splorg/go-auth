package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/splorg/go-auth/internal/util"
)

func Authenticate(c *fiber.Ctx) error {
    jwtCookie := c.Cookies("jwt")

    if jwtCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    token, err := jwt.ParseWithClaims(jwtCookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(util.JwtSecret), nil
    })
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    if !token.Valid {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    claims, ok := token.Claims.(*jwt.StandardClaims)
    if !ok || claims.ExpiresAt < time.Now().Unix() {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    return c.Next()
}