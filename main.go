package main

import (

	// "fmt"
	"context"
	"fmt"
	"log"
	"time"

	db "github.com/SeongUgKim/simplebank/db/postgres"
	accountentity "github.com/SeongUgKim/simplebank/entity/account"
	accountrepository "github.com/SeongUgKim/simplebank/repository/account"
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()

	postgresDB, err := db.New("postgres", "postgres://postgres:mysecretpassword@localhost:5433/simple_bank?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo, err := accountrepository.New(accountrepository.Params{
		DB: postgresDB,
	})
	if err != nil {
		log.Fatal(err)
	}

	accountUUID, err := uuid.NewV4()
	if err != nil {
		log.Fatal(err)
	}

	account := accountentity.Account{
		UUID:      accountUUID.String(),
		Owner:     "Seong Ug Kim",
		AmountE5:  10000000,
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	insertedAccount, err := repo.Insert(ctx, account)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s, %s, %d, %s\n", insertedAccount.UUID, insertedAccount.Owner, insertedAccount.AmountE5, insertedAccount.Currency)
	updateAccount := accountentity.Account{
		UUID:      insertedAccount.UUID,
		Owner:     insertedAccount.Owner,
		AmountE5:  30000000,
		Currency:  "USD",
		CreatedAt: time.Now(),
	}

	updatedAccount, err := repo.Update(ctx, updateAccount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s, %s, %d, %s\n", updatedAccount.UUID, updatedAccount.Owner, updatedAccount.AmountE5, updatedAccount.Currency)

	if err := repo.Delete(ctx, updateAccount.UUID); err != nil {
		log.Fatal(err)
	}
}
