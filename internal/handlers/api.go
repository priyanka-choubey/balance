package handlers

import (
	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/priyanka-choubey/balance/internal/middleware"
)

func Handler(r *chi.Mux) {
	// Global middleware
	r.Use(chimiddle.StripSlashes)

	r.Route("/account", func(router chi.Router) {

		// Middleware for /account route
		router.Use(middleware.Authorization)

		router.Get("/coins", GetBalance)
		router.Delete("/delete", DeleteUser)
		router.Put("/updatetoken", UpdateUserToken)
		// TODO: router.Put("/update", UpdateAccount)
	})

	r.Route("/user", func(router chi.Router) {

		router.Post("/create", CreateUser)
	})
}
