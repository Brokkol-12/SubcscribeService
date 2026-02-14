package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"subscription-service/configs"

	"github.com/jackc/pgx/v4"
)

func main() {
	cfg := configs.LoadConfig()
	conn, err := pgx.Connect(context.Background(), cfg.Db.Dsn)
	if err != nil {
		log.Fatal("Error to set connection to BD", err)
	}
	defer conn.Close(context.Background())

	sqlBytes, err := os.ReadFile("internal/db/subs_table.sql")
	if err != nil {
		log.Fatal("Read SQL file error ", err)
	}
	_, err = conn.Exec(context.Background(), string(sqlBytes))
	if err != nil {
		log.Fatal("Migration ERROR ", err)
	}
	fmt.Println("Migration success")
}
