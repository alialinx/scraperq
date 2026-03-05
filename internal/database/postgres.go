package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type DB struct {
	Pool *pgxpool.Pool
}

func NewDB(user string, password string, host string, port string, dbname string) (*DB, error) {

	db := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	pool, err := pgxpool.New(ctx, db)

	if err != nil {
		return nil, fmt.Errorf("db connection error: %w", err)
	}

	err = pool.Ping(ctx)

	if err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}

	return &DB{Pool: pool}, nil

}

func (db *DB) Close() {
	db.Pool.Close()
}
