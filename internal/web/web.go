package web

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error) {
	var clientErr *ClientError
	if errors.As(err, &clientErr) {
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, ClientError{
			Code: clientErr.Code,
			Err:  err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	log.Error().
		Stack().
		Err(err).
		Msg("web: Internal server error")
}
