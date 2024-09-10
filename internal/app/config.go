package app

import (
	"fmt"
)

const urlPrefix = "/api"

type (
	Options struct {
		Addr string
	}

	path struct {
		ping string
	}
	config struct {
		addr string
		path path
	}
)

func NewConfig(opts Options) config {
	return config{
		addr: opts.Addr,
		path: path{
			ping: fmt.Sprintf("GET %s/ping", urlPrefix),
		},
	}
}
