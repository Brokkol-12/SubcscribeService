package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPool(conn string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), conn)
}
