package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
)

func MakeRouter(_ context.Context, handler *BtcPrice) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(30 * time.Second)) //nolint:gomnd

	router.Route("/api", func(r chi.Router) {
		r.Get("/rate", handler.handleRate)
		r.Post("/subscribe", handler.handleSubscribe)
		r.Get("/sendEmails", handler.handleSendEmails)
	})

	return router
}
