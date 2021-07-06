package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	//ID         int `gorm:"private_key" json:"id"`
	CreatedOn  int            `gorm:"autoCreateTime" json:"created_on"`
	ModifiedOn int            `gorm:"autoUpdateTime:milli" json:"modified_on"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func Setup(test bool) {
	var err error

	if !test {
		err = godotenv.Load()
		if err != nil {
			log.Fatalf("models.Setup err: %v", err)
		}
	} else {
		err = godotenv.Load("../.env")
		if err != nil {
			log.Fatalf("models.Setup err: %v", err)
		}
	}

	fmt.Println(os.Getenv("DATABASE_DSN"))

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DATABASE_DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db.AutoMigrate(&Shoppinglist{})
	db.AutoMigrate(&Auth{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
