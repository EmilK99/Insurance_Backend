package sc

import (
	"context"
	"github.com/mr-tron/base58"
	"github.com/portto/solana-go-sdk/client"
	"github.com/portto/solana-go-sdk/common"
	"github.com/portto/solana-go-sdk/pkg/bincode"
	"github.com/portto/solana-go-sdk/program/sysprog"
	"github.com/portto/solana-go-sdk/types"
	"log"
)

var key = []byte{239, 135, 109, 127, 74, 161, 217, 168, 151, 232, 108, 167, 47, 189, 243, 246, 126, 215, 7, 209, 223, 231, 174, 124, 129, 82, 222, 251, 212, 186, 137, 242, 140, 230, 149, 19, 121, 132, 205, 249, 133, 114, 200, 173, 189, 139, 120, 79, 87, 207, 112, 93, 201, 147, 1, 136, 92, 172, 123, 165, 67, 116, 60, 254}

type Client struct {
	Client    *client.Client
	Account   types.Account
	ProgramID common.PublicKey
}

type Instruction uint8

const (
	InstructionSave Instruction = iota
	InstructionClose
	InsuranceAccountSize = 6
)

func NewClient(programID common.PublicKey, privateKey, endpoint string) (*Client, error) {
	acc, err := types.AccountFromBase58(privateKey)
	if err != nil {
		return nil, err
	}

	c := client.NewClient(endpoint)

	return &Client{Client: c, Account: acc, ProgramID: programID}, nil
}

func (cl *Client) CreateInsuranceContract(ctx context.Context, contractID int) (string, error) {
	contractAccount := types.NewAccount()
	log.Println("contract account:", contractAccount.PublicKey.ToBase58())
	log.Println(base58.Encode(contractAccount.PrivateKey))

	rentExemptionBalance, err := cl.Client.GetMinimumBalanceForRentExemption(ctx, InsuranceAccountSize)
	if err != nil {
		return "", err
	}

	res, err := cl.Client.GetRecentBlockhash(ctx)
	if err != nil {
		return "", err
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			sysprog.CreateAccount(
				cl.Account.PublicKey,
				contractAccount.PublicKey,
				cl.ProgramID,
				rentExemptionBalance,
				InsuranceAccountSize,
			),
			cl.saveContract(
				cl.Account.PublicKey,
				contractAccount.PublicKey,
				uint32(contractID),
			),
		},
		Signers:         []types.Account{cl.Account, contractAccount},
		FeePayer:        cl.Account.PublicKey,
		RecentBlockHash: res.Blockhash,
	})
	if err != nil {
		return "", err
	}

	txhash, err := cl.Client.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return "", err
	}

	log.Println("txhash:", txhash)

	return base58.Encode(contractAccount.PrivateKey), nil
}

func (cl *Client) CloseInsuranceContract(ctx context.Context, accountKey string) error {
	contractAccount, err := types.AccountFromBase58(accountKey)
	if err != nil {
		return err
	}

	blockHash, err := cl.Client.GetRecentBlockhash(ctx)
	if err != nil {
		return err
	}

	rawTx, err := types.CreateRawTransaction(types.CreateRawTransactionParam{
		Instructions: []types.Instruction{
			cl.closeContract(
				cl.Account.PublicKey,
				contractAccount.PublicKey,
			),
		},
		Signers:         []types.Account{cl.Account},
		FeePayer:        cl.Account.PublicKey,
		RecentBlockHash: blockHash.Blockhash,
	})

	if err != nil {
		return err
	}

	txhash, err := cl.Client.SendRawTransaction(ctx, rawTx)
	if err != nil {
		return err
	}

	log.Println("txhash:", txhash)

	return nil
}

// SaveContract init a contract with specified id
func (c *Client) saveContract(ownerPublickey, accountPublicKey common.PublicKey, contractID uint32) types.Instruction {
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

// closeContract closed contract with specified id
func (c *Client) closeContract(ownerPublickey, accountPublicKey common.PublicKey) types.Instruction {
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
