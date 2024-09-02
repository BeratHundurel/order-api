package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

type App struct {
	router       http.Handler
	rdb          *redis.Client
	cfg          Config
	currencyConn *grpc.ClientConn
}

func New(config Config) *App {
	conn, err := grpc.NewClient(config.gRPCPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to currency service: %v", err)
	}
	defer conn.Close()

	app := &App{
		rdb: redis.NewClient(&redis.Options{
			Addr: config.RedisAddress,
		}),
		cfg:          config,
		currencyConn: conn,
	}

	app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.ServerPort),
		Handler: a.router,
	}

	err := a.rdb.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	defer func() {
		if err := a.rdb.Close(); err != nil {
			fmt.Println("failed to close redis", err)
		}
	}()

	fmt.Println("Starting server")

	ch := make(chan error, 1)

	go func() {
		err = server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed to start server: %w", err)
		}
		close(ch)
	}()

	select {
	case err = <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		return server.Shutdown(timeout)
	}
}
