package store

import (
	"context"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/types"
)

func (s *Store) SaveContractAccount(ctx context.Context, contractID int, account *types.Account) error {
	_, err := s.Conn.Exec(ctx, "UPDATE contracts SET sc_account = $2, sc_key = $3 WHERE id = $1",
		contractID, account.PublicKey.ToBase58(), base58.Encode(account.PrivateKey))
	if err != nil {
		return err
	}

	return nil
}
