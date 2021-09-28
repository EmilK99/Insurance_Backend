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

func (s Store) FindContracts(user_id string) ([]Contract, error) {

	return nil, nil
}
