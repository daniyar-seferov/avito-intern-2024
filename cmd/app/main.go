package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"avito/tender/internal/app"
)

func main() {
	initOpts()
	initMigration(opts.DBConnStr)
	service, err := app.NewApp(app.NewConfig(opts))
	if err != nil {
		log.Fatal("{FATAL} ", err)
	}

	go func() {
		log.Printf("Starting server on %s\n", opts.Addr)
		err = service.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server closed\n")
		}
		if err != nil {
			log.Fatalf("error starting server: %s\n", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop
	fmt.Println("Shutting down...")
	if err := service.Close(); err != nil {
		log.Printf("Error closing service: %v", err)
	}
}
