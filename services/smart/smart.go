package smart

import (
	votingpb "github.com/egsam98/voting/proto"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"google.golang.org/protobuf/proto"
)

const (
	funcRegisterVote = "RegisterVote"
	funcFindVote     = "FindVote"
)

type Contract struct {
	contractapi.Contract
}

func (*Contract) RegisterVote(ctx contractapi.TransactionContextInterface, votePbStr string) error {
	votePb := []byte(votePbStr)

	vote := &votingpb.Vote{}
	if err := proto.Unmarshal(votePb, vote); err != nil {
		return err
	}

	key := getVoteKey(vote.CandidateId, vote.Voter.GetPassport())

	b, err := ctx.GetStub().GetState(key)
	if err != nil {
		return err
	}

	if b != nil {
		vote.Status = votingpb.Vote_FAIL
		reason := "vote by user with this passport already exists"
		vote.FailReason = &reason

		if votePb, err = proto.Marshal(vote); err != nil {
			return err
		}
	}

	return ctx.GetStub().PutState(key, votePb)
}

func (*Contract) FindVote(ctx contractapi.TransactionContextInterface, candidateID int64, voterPassport string) (string, error) {
	key := getVoteKey(candidateID, voterPassport)
	b, err := ctx.GetStub().GetState(key)
	return string(b), err
}
