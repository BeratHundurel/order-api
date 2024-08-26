package application

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// router.Route("/orders", a.loadOrderRoutes)

	a.router = router
}

// func (a *App) loadOrderRoutes(router chi.Router) {
// 	orderHandler := &order.OrderHandler{
// 		Repo: &order.RedisRepo{
// 			Client: a.rdb,
// 		},
// 	}

// 	router.Post("/", orderHandler.Create)
// 	router.Get("/", orderHandler.List)
// 	router.Get("/{id}", orderHandler.GetByID)
// 	router.Put("/{id}", orderHandler.UpdateByID)
// 	router.Delete("/{id}", orderHandler.DeleteByID)
// }
