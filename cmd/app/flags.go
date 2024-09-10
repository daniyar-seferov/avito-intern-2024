package main

import (
	"os"

	"avito/tender/internal/app"
)

const (
	// defaultAddr = ":8080"

	envServerAddress = "SERVER_ADDRESS"
	envPostgresConn  = "POSTGRES_CONN"
)

var opts = app.Options{}

func initOpts() {
	// flag.StringVar(&opts.Addr, "addr", defaultAddr, fmt.Sprintf("server address, default: %q", defaultAddr))
	// flag.Parse()

	opts.Addr = os.Getenv(envServerAddress)
	opts.DBConnStr = os.Getenv(envPostgresConn)
}
