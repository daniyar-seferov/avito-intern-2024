package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func initMigration(dbConnStr string) {
	var db *sql.DB
	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	fmt.Println("Starting migrations")
	if err := goose.Up(db, "migrations"); err != nil {
		log.Fatalf("error migration up: %s\n", err)
	}
}
