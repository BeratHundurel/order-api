package application

import (
	"net/http"

	"github.com/BeratHundurel/order-api/authentication-api/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/auth", a.loadAuthRoutes)

	a.router = router
}

func (a *App) loadAuthRoutes(router chi.Router) {
	authHandler := &auth.Handler{
		Repo: &auth.TursoRepo{
			DB: a.db,
		},
	}

	router.Post("/", authHandler.Register)
	router.Get("/{id}", authHandler.GetByID)
}
