package main

import (
	"errors"
	"log"
	"net/http"

	"avito/tender/internal/app"
)

func main() {
	initOpts()
	service, err := app.NewApp(app.NewConfig(opts))
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}

	log.Printf("Starting server on %s\n", opts.Addr)
	err = service.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		log.Fatal("server closed\n")
	}
	if err != nil {
		log.Fatalf("error starting server: %s\n", err)
	}
}
