package sc

import (
	"context"
	"github.com/portto/solana-go-sdk/types"
)

type Client struct {
	Account types.Account
}

func NewClient(key []byte) (*Client, error) {
	acc, err := types.AccountFromBytes(key)
	if err != nil {
		return nil, err
	}

	return &Client{Account: acc}, nil
}

func (client *Client) CreateSmartContract(ctx context.Context, id int) string {

	//TODO: implement smart contract interaction

	return ""
}
