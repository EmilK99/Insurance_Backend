package sc

import (
	"context"
	"fmt"
	client2 "github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/client/rpc"
	"github.com/portto/solana-go-sdk/types"
	"log"
)

var key = []byte{239, 135, 109, 127, 74, 161, 217, 168, 151, 232, 108, 167, 47, 189, 243, 246, 126, 215, 7, 209, 223, 231, 174, 124, 129, 82, 222, 251, 212, 186, 137, 242, 140, 230, 149, 19, 121, 132, 205, 249, 133, 114, 200, 173, 189, 139, 120, 79, 87, 207, 112, 93, 201, 147, 1, 136, 92, 172, 123, 165, 67, 116, 60, 254}

type Client struct {
	Account types.Account
}

func NewClient() (*Client, error) {
	acc, err := types.AccountFromBytes(key)
	if err != nil {
		return nil, err
	}

	return &Client{Account: acc}, nil
}

func (cl *Client) CreateSmartContract(ctx context.Context, id int) string {
	c := client2.NewClient(rpc.DevnetRPCEndpoint)
	balance, err := c.GetBalance(
		ctx,
		cl.Account.PublicKey.ToBase58(),
	)
	if err != nil {
		log.Fatalln("get balance error", err)
	}
	fmt.Println(balance)

	//TODO: implement smart contract interaction

	return ""
}
