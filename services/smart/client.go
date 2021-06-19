package smart

import (
	votingpb "github.com/egsam98/voting/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	contract *gateway.Contract
}

func NewClient(contract *gateway.Contract) *Client {
	return &Client{contract: contract}
}

func (c *Client) RegisterVote(vote *votingpb.Vote) error {
	b, err := proto.Marshal(vote)
	if err != nil {
		return err
	}

	if _, err = c.contract.SubmitTransaction("RegisterVote", string(b)); err != nil {
		return err
	}

	return nil
}
