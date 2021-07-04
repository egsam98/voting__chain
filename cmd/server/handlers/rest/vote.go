package rest

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/egsam98/voting/chain/internal/web"
	"github.com/egsam98/voting/chain/services/smart"
)

type voteController struct {
	client *smart.Client
}

func newVoteController(client *smart.Client) *voteController {
	return &voteController{client: client}
}

func (vc *voteController) FindVoteByVoterPassport(w http.ResponseWriter, r *http.Request) {
	candidateID, err := strconv.ParseInt(chi.URLParam(r, "candidate_id"), 10, 64)
	if err != nil {
		web.RespondError(w, r, web.WrapWithError(smart.ErrInvalidInput, err))
		return
	}

	passport := chi.URLParam(r, "passport")

	vote, err := vc.client.FindVote(candidateID, passport)
	if err != nil {
		web.RespondError(w, r, err)
		return
	}

	render.JSON(w, r, vote)
}
