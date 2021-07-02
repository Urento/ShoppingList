package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	ID         int `gorm:"private_key" json:"id"`
	CreatedOn  int `gorm:"autoCreateTime" json:"created_on"`
	ModifiedOn int `gorm:"autoUpdateTime:milli" json:"modified_on"`
}

func Setup() {
	var err error

	err = godotenv.Load()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("POSTGRES_DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	//TODO: Automigrate (https://gorm.io/docs/index.html)

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
