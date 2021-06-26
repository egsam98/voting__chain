package rest

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

func API(admin sarama.ClusterAdmin, producer sarama.SyncProducer) http.Handler {
	hc := newHealthController(admin, producer)

	mux := chi.NewMux()
	mux.Use(
		middleware.RequestLogger(&middleware.DefaultLogFormatter{
			Logger: &log.Logger,
		}),
		middleware.Recoverer,
	)

	mux.Route("/health", func(r chi.Router) {
		r.Get("/readiness", hc.Readiness)
	})

	api := chi.NewRouter()
	mux.Mount("/api", api)

	return mux
}
