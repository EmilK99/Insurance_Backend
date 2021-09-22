package store

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	Pool *pgxpool.Pool
	Conn *pgxpool.Conn
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{Pool: pool}
}

func (s Store) FindContracts(user_id string) ([]Contract, error) {

	return nil, nil
}
