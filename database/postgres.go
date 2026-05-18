package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{DB: db}
}

func NewPostgresConnection() (*pgxpool.Pool, error) {

	dsn := "postgres://postgres:pass1234@localhost:5432/Lifemetrics"

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Database connection established successfully ✅")

	return pool, nil
}
