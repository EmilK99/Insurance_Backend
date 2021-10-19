package store

import "context"

func (s *Store) SaveContractAccount(ctx context.Context, contractID int, accountKey string) error {
	_, err := s.Conn.Exec(ctx, "UPDATE contracts SET sc_account = $2 WHERE id = $1", contractID, accountKey)
	if err != nil {
		return err
	}

	return nil
}
