package smart

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Contract struct {
	contractapi.Contract
}

func (c *Contract) RegisterVote(ctx contractapi.TransactionContextInterface, key, value string) error {
	return nil
}
