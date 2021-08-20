package app

import (
	"flight_app/app/contract"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Store struct {
	pool *pgxpool.Pool
}

func NewStore(pool *pgxpool.Pool) *Store {
	return &Store{pool: pool}
}

func (s Store) FindContracts(user_id string) ([]contract.Contract, error) {

	return nil, nil
}
