package store

import (
	"github.com/jackc/pgx/v5"
)

type Store struct {
	DB *DB
	// Redis
}

func New(conn *pgx.Conn) *Store {
	return &Store{
		DB: NewDB(conn),
	}
}
