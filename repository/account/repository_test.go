package account

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	db "github.com/SeongUgKim/simplebank/db/postgres"
	entity "github.com/SeongUgKim/simplebank/entity/account"
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testRepo Repository

const (
	dbDriver = "postgres"
	dbSource = "postgres://postgres:mysecretpassword@localhost:5433/simple_bank?sslmode=disable"
)

var (
	date  = time.Date(2024, time.May, 10, 0, 0, 0, 0, time.UTC)
	uuids = []string{
		uuid.Must(uuid.NewV4()).String(),
		uuid.Must(uuid.NewV4()).String(),
		uuid.Must(uuid.NewV4()).String(),
	}
	accounts = []entity.Account{
		{UUID: uuids[0], Owner: "Seong Ug Kim", AmountE5: 100000000, Currency: "USD", CreatedAt: date},
		{UUID: uuids[1], Owner: "James Kosher", AmountE5: 1500000000, Currency: "USD", CreatedAt: date},
		{UUID: uuids[2], Owner: "Travis Scott", AmountE5: 200000000, Currency: "USD", CreatedAt: date},
	}
)

func TestMain(m *testing.M) {
	postgresDB, err := db.New(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testRepo, err = New(Params{DB: postgresDB})
	if err != nil {
		log.Fatal("cannot create test db:", err)
	}

	os.Exit(m.Run())
}

func TestInsert(t *testing.T) {

	testCases := map[string]struct {
		account entity.Account
		res     entity.Account
		err     string
	}{
		"success": {
			account: accounts[0],
			res:     accounts[0],
			err:     "",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			res, err := testRepo.Insert(ctx, tc.account)
			res.CreatedAt = tc.res.CreatedAt
			assert.Equal(t, tc.res, res)
			if tc.err != "" {
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
