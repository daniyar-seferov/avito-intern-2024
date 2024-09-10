package main

import (
	"flag"
	"fmt"

	"avito/tender/internal/app"
)

const (
	defaultAddr = ":8080"
)

var opts = app.Options{}

func initOpts() {
	flag.StringVar(&opts.Addr, "addr", defaultAddr, fmt.Sprintf("server address, default: %q", defaultAddr))
	flag.Parse()
}
