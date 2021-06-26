package smart

import (
	votingpb "github.com/egsam98/voting/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/multi"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/pkg/errors"
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
		return errors.Wrapf(err, "failed to marhal %T to protobuf", vote)
	}

	_, err = c.contract.SubmitTransaction(funcRegisterVote, string(b))
	if err == nil {
		return nil
	}

	var multiErr multi.Errors
	if !errors.As(err, &multiErr) {
		return errors.Errorf("failed to decode contract error, expected %T", multiErr)
	}

	for _, err := range multiErr {
		var s *status.Status
		if !errors.As(err, &s) {
			return errors.Errorf("failed to decode one of multi errors, expected %T", s)
		}

		if s.Code == 500 {
			return errors.Wrap(ErrInvalidVote, s.Message)
		}
	}

	return errors.Wrapf(err, "failed to invoke %q", funcRegisterVote)
}
