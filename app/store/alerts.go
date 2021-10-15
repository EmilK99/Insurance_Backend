package store

import (
	"context"
)

func (s *Store) UpdateContractsByAlert(ctx context.Context, flightID, status string, flightDate int64) error {

	_, err := s.Conn.Exec(ctx, "UPDATE contracts SET status = $3 WHERE flight_number = $1 AND "+
		"flight_date = $2 AND status = 'waiting'", flightID, flightDate, status)

	return err
}
