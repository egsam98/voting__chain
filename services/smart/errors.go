package smart

import (
	"errors"

	"github.com/egsam98/voting/chain/internal/web"
)

var ErrInvalidVote = errors.New("vote is invalid")

var (
	ErrInvalidInput = &web.ClientError{
		Code: 1,
		Err:  "invalid input",
	}
	ErrVoteNotFound = &web.ClientError{
		Code: 2,
		Err:  "vote is not found",
	}
)
