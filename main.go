package main

import (
	"log"

	"github.com/SeongUgKim/simplebank/config"
	accountcontroller "github.com/SeongUgKim/simplebank/controller/account"
	db "github.com/SeongUgKim/simplebank/db/postgres"
	accounthandler "github.com/SeongUgKim/simplebank/handler/account"
	accountrepository "github.com/SeongUgKim/simplebank/repository/account"
	_ "github.com/lib/pq"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}

	postgresDB, err := db.New(config.DBDriver, config.DBSource)
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

	if err := handler.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
