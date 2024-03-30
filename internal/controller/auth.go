package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/splorg/go-auth/internal/database"
	"github.com/splorg/go-auth/internal/dto"
	"github.com/splorg/go-auth/internal/model"
	"github.com/splorg/go-auth/internal/util"
	"github.com/splorg/go-auth/internal/validator"
)

func HelloWorld(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Hello, world!"})
}

func Register(c *fiber.Ctx) error {
	var req dto.RegisterDTO

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	password, err := util.HashPassword([]byte(req.Password))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to encrypt password"})
	}

	user := model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: password,
	}

	result := database.DB.Omit("password").Create(&user)
	if result.RowsAffected == 0 {
		c.SendStatus(400)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Credentials already in use."})
	}

	token, err := util.GenerateJWT(&user, util.JwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT"})
	}

	cookie := util.GenerateCookie(token)

	c.Cookie(cookie)

	return c.Status(fiber.StatusCreated).JSON(user)
}

func Login(c *fiber.Ctx) error {
	var req dto.LoginDTO

	if err := c.BodyParser(&req); err != nil {
		return err
	}

	if err := validator.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var user model.User

	database.DB.Where("email = ?", req.Email).First(&user)

	if user.ID == "" {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := util.ComparePassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := util.GenerateJWT(&user, util.JwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate JWT"})
	}

	cookie := util.GenerateCookie(token)

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(user)
}

func ProtectedRoute(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"success": true, "token": c.Cookies("jwt")})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(util.JwtSecret), nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims := token.Claims.(*jwt.StandardClaims)

	var user model.User

	database.DB.Where("id = ?", claims.Subject).First(&user)

	return c.JSON(user)
}
