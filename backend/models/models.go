package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	utils "github.com/urento/shoppinglist/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Model struct {
	CreatedOn  int            `gorm:"autoCreateTime" json:"created_on"`
	ModifiedOn int            `gorm:"autoUpdateTime:milli" json:"modified_on"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
}

func Setup() {
	var err error

	if utils.PROD {
		err = godotenv.Load()
	} else if utils.GITHUB_TESTING {
		err = nil
	} else {
		err = godotenv.Load("../.env")
	}

	if err != nil {
		log.Fatal(err)
	}

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
