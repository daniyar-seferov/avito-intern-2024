package app

import (
	"context"
	"log"
	"net/http"
	"time"

	appHttp "avito/tender/internal/app/http"

	pgxv5 "github.com/jackc/pgx/v5"
)

type (
	mux interface {
		Handle(pattern string, handler http.Handler)
	}
	server interface {
		ListenAndServe() error
		Close() error
	}

	App struct {
		config config
		mux    mux
		server server
	}
)

func NewApp(config config) (*App, error) {
	var mux = http.NewServeMux()

	ctx := context.Background()
	conn, err := pgxv5.Connect(ctx, config.dbConnStr)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &App{
		config: config,
		mux:    mux,
		server: &http.Server{Addr: config.addr, Handler: wrapLogger(mux)},
	}, nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.ping, appHttp.NewPingHandler())

	return a.server.ListenAndServe()
}

func (a *App) Close() error {
	return a.server.Close()
}

func wrapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
