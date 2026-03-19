package main

import (
	"log"
	"net/http"

	"github.com/ahmedkltn/tenant-api/internal/auth"
	"github.com/ahmedkltn/tenant-api/internal/middleware"
	"github.com/ahmedkltn/tenant-api/internal/projects"
	"github.com/ahmedkltn/tenant-api/internal/seed"
	"github.com/ahmedkltn/tenant-api/internal/store"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// setup
	s := store.New()
	seed.Load(s)
	h := projects.New(s)

	// router
	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	// public
	r.Post("/auth/login", auth.NewLoginHandler(s))

	// protected
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		r.Route("/tenants/{tenantId}/projects", func(r chi.Router) {
			r.Use(middleware.TenantGuard)

			r.Get("/", h.List)
			r.With(middleware.RoleCheck(store.RoleAdmin)).Post("/", h.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Get)
				r.With(middleware.RoleCheck(store.RoleAdmin)).Delete("/", h.Delete)

			})
		})
	})

	log.Println("server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
