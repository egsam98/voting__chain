package smart

import (
	"fmt"
)

func getVoteKey(candidateID int64, voterPassport string) string {
	return fmt.Sprintf("vote (candidate=%d passport=%s)", candidateID, voterPassport)
}
