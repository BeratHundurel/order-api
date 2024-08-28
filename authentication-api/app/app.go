package application

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type App struct {
	router http.Handler
	cfg    Config
	db     *gorm.DB
}

func New(config Config) *App {
	if err := InitializeDB(); err != nil {
		fmt.Println("Failed to initialize the database:", err)
		os.Exit(1)
	}

	app := &App{
		cfg: config,
		db:  GetDB(),
	}
	app.loadRoutes()

	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.ServerPort),
		Handler: a.router,
	}

	fmt.Println("Starting server")

	MigrateDB()

	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
