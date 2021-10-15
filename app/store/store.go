package store

import (
	"github.com/jackc/pgx/v4"
)

type Store struct {
	Conn *pgx.Conn
}

func NewStore(conn *pgx.Conn) *Store {
	return &Store{Conn: conn}
}
