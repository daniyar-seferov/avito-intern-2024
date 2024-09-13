package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	appHttp "avito/tender/internal/app/http"
	"avito/tender/internal/domain"
	tenders_new "avito/tender/internal/handlers/tenders/new"
	db_pgx_repo "avito/tender/internal/repository/pgx"

	pgxv5 "github.com/jackc/pgx/v5"
)

type (
	mux interface {
		Handle(pattern string, handler http.Handler)
	}
	server interface {
		ListenAndServe() error
		Shutdown(ctx context.Context) error
	}
	tenderStorage interface {
		AddTender(ctx context.Context, item domain.TenderAddDTO) (string, error)
		GetUserOrganizationId(ctx context.Context, username string) (string, string, error)
		GetTender(ctx context.Context, tenderId string) (domain.TenderAddDTO, error)
	}

	App struct {
		config  config
		mux     mux
		server  server
		dbConn  *pgxv5.Conn
		storage tenderStorage
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
		config:  config,
		mux:     mux,
		server:  &http.Server{Addr: config.addr, Handler: wrapLogger(mux)},
		dbConn:  conn,
		storage: db_pgx_repo.NewRepo(conn),
	}, nil
}

func (a *App) ListenAndServe() error {
	a.mux.Handle(a.config.path.ping, appHttp.NewPingHandler())
	a.mux.Handle(a.config.path.tendersAdd, appHttp.NewTendersAddHandler(tenders_new.New(a.storage), a.config.path.tendersAdd))

	return a.server.ListenAndServe()
}

func (a *App) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown error: %v", err)
	}

	if err := a.dbConn.Close(ctx); err != nil {
		return fmt.Errorf("failed to close the database connection: %v", err)
	}

	return nil
}

func wrapLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}
