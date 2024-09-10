package main

import (
	"avito/tender/internal/app"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const (
	envServerAddress = "SERVER_ADDRESS"
	envPostgresConn  = "POSTGRES_CONN"
	envMode          = "MODE"
)

var opts = app.Options{}

func initOpts() {
	mode := os.Getenv(envMode)
	if mode != "dev" {
		log.Println("Prod environment")
		err := godotenv.Load("./build/prod/.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		log.Println("Dev environment")
	}

	opts.Addr = os.Getenv(envServerAddress)
	opts.DBConnStr = os.Getenv(envPostgresConn)
}
