package util

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	JwtSecret = os.Getenv("JWT_SECRET")
	DatabaseURL = os.Getenv("DATABASE_URL")
)