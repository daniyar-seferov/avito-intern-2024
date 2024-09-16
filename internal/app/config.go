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
		ping, tendersAdd, tendersList, tendersMy, tendersStatus, tendersChangeStatus string
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
			ping:                fmt.Sprintf("GET %s/ping", urlPrefix),
			tendersAdd:          fmt.Sprintf("POST %s/tenders/new", urlPrefix),
			tendersList:         fmt.Sprintf("GET %s/tenders", urlPrefix),
			tendersMy:           fmt.Sprintf("GET %s/tenders/my", urlPrefix),
			tendersStatus:       fmt.Sprintf("GET %s/tenders/{tenderId}/status", urlPrefix),
			tendersChangeStatus: fmt.Sprintf("PUT %s/tenders/{tenderId}/status", urlPrefix),
		},
	}
}
