package main

import (
	"log"

	accountcontroller "github.com/SeongUgKim/simplebank/controller/account"
	db "github.com/SeongUgKim/simplebank/db/postgres"
	accounthandler "github.com/SeongUgKim/simplebank/handler/account"
	accountrepository "github.com/SeongUgKim/simplebank/repository/account"
	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgres://postgres:mysecretpassword@localhost:5433/simple_bank?sslmode=disable"
	serverAddress = "0.0.0.0:3000"
)

func main() {
	postgresDB, err := db.New(dbDriver, dbSource)
	if err != nil {
		log.Fatal(err)
	}

	repo, err := accountrepository.New(accountrepository.Params{
		DB: postgresDB,
	})
	if err != nil {
		log.Fatal(err)
	}

	controller, err := accountcontroller.New(accountcontroller.Params{
		Repository: repo,
	})
	if err != nil {
		log.Fatal(err)
	}

	handler, err := accounthandler.New(accounthandler.Params{
		Controller: controller,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := handler.Start(serverAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

}
