package app

import (
	"fmt"
)

const urlPrefix = "/api"

type (
	Options struct {
		Addr      string
		DBConnStr string
	}

	path struct {
		ping, tendersAdd string
	}
	config struct {
		addr      string
		dbConnStr string
		path      path
	}
)

func NewConfig(opts Options) config {
	return config{
		addr:      opts.Addr,
		dbConnStr: opts.DBConnStr,
		path: path{
			ping:       fmt.Sprintf("GET %s/ping", urlPrefix),
			tendersAdd: fmt.Sprintf("POST %s/tenders/new", urlPrefix),
		},
	}
}
