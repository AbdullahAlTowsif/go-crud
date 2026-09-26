package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

var Db *pgx.Conn

func ConnectDb() {
	var err error
	connStr := os.Getenv("DB_URL")

	Db, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connected Successfully!")
}
