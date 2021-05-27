package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	util "github.com/urento/shoppinglist/pkg"
)

var DB *pgxpool.Pool

func Connect() {
	var err error
	//check if prod or local eviorment
	if !util.PROD {
		err = godotenv.Load("../.env")
	} else {
		err = godotenv.Load()
	}
	if err != nil {
		panic(err)
	}

	// connect to db
	poolConfig, err := pgxpool.ParseConfig((os.Getenv("DATABASE_URL")))
	if err != nil {
		panic(err)
	}

	DB, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	defer DB.Close()
}
