package smart

import (
	votingpb "github.com/egsam98/voting/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"google.golang.org/protobuf/proto"
)

type Contract struct {
	contractapi.Contract
}

func (*Contract) RegisterVote(ctx contractapi.TransactionContextInterface, voteProto string) error {
	vote := &votingpb.Vote{}
	if err := proto.Unmarshal([]byte(voteProto), vote); err != nil {
		return err
	}

	return nil
}
