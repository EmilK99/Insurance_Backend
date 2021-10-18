package sc

import (
	"context"
	"fmt"
	client2 "github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/client/rpc"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/bincode"
	"github.com/portto/solana-go-sdk/types"
	"log"
)

var key = []byte{239, 135, 109, 127, 74, 161, 217, 168, 151, 232, 108, 167, 47, 189, 243, 246, 126, 215, 7, 209, 223, 231, 174, 124, 129, 82, 222, 251, 212, 186, 137, 242, 140, 230, 149, 19, 121, 132, 205, 249, 133, 114, 200, 173, 189, 139, 120, 79, 87, 207, 112, 93, 201, 147, 1, 136, 92, 172, 123, 165, 67, 116, 60, 254}

type Client struct {
	Account   types.Account
	ProgramID common.PublicKey
}

type Instruction uint8

const (
	InstructionSave Instruction = iota
	InstructionClose
)

func NewClient(programID common.PublicKey) (*Client, error) {
	acc, err := types.AccountFromBytes(key)
	if err != nil {
		return nil, err
	}

	return &Client{Account: acc, ProgramID: programID}, nil
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

// SaveContract init a contract with specified id
func (c *Client) SaveContract(ownerPublickey, accountPublicKey common.PublicKey, contractID uint32) types.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction         Instruction
		InsuranceContractID uint32
	}{
		Instruction:         InstructionSave,
		InsuranceContractID: contractID,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: ownerPublickey, IsSigner: true, IsWritable: false},
		{PubKey: accountPublicKey, IsSigner: false, IsWritable: true},
		{PubKey: common.SysVarRentPubkey, IsSigner: false, IsWritable: false},
	}
	return types.Instruction{
		ProgramID: c.ProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}

// CloseContract closed contract with specified id
func (c *Client) CloseContract(ownerPublickey, accountPublicKey common.PublicKey) types.Instruction {
	data, err := bincode.SerializeData(struct {
		Instruction Instruction
	}{
		Instruction: InstructionClose,
	})
	if err != nil {
		panic(err)
	}

	accounts := []types.AccountMeta{
		{PubKey: ownerPublickey, IsSigner: true, IsWritable: false},
		{PubKey: accountPublicKey, IsSigner: false, IsWritable: true},
	}

	return types.Instruction{
		ProgramID: c.ProgramID,
		Accounts:  accounts,
		Data:      data,
	}
}
