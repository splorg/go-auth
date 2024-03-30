package database

import (
	"log"

	"github.com/splorg/go-auth/internal/model"
	"github.com/splorg/go-auth/internal/util"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open(util.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

  DB = db

	db.AutoMigrate(&model.User{})
}
