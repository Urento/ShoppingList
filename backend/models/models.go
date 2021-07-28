package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	utils "github.com/urento/shoppinglist/pkg"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	} else if utils.LOCAL_TESTING {
		err = godotenv.Load("../.env")
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal(err)
	}

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)

	db, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  os.Getenv("DATABASE_DSN"),
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalf("Error while connecting to database: %s", err)
	}

	db.AutoMigrate(&Shoppinglist{})
	db.AutoMigrate(&Auth{})
	db.AutoMigrate(&ResetPassword{})

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error while connecting to database: %s", err)
	}

	sqlDB.SetMaxIdleConns(1000)
	sqlDB.SetMaxOpenConns(10000)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
