package account

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/SeongUgKim/simplebank/config"
	db "github.com/SeongUgKim/simplebank/db/postgres"
	entity "github.com/SeongUgKim/simplebank/entity/account"
	"github.com/gofrs/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testRepo Repository

func TestMain(m *testing.M) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	postgresDB, err := db.New(config.DBDriver, config.DBSource)
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

func TestFetch(t *testing.T) {
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

	testCases := map[string]struct {
		insertedAccounts []entity.Account
		accountUUID      string
		res              entity.Account
		err              string
	}{
		"success": {
			insertedAccounts: accounts,
			accountUUID:      uuids[0],
			res:              accounts[0],
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()

			for _, account := range tc.insertedAccounts {
				acc, err := testRepo.Insert(ctx, account)
				assert.NoError(t, err)
				assert.NotNil(t, acc)
			}

			res, err := testRepo.Fetch(ctx, tc.accountUUID)
			res.CreatedAt = tc.res.CreatedAt
			assert.Equal(t, tc.res, res)
			if err != nil {
				assert.EqualError(t, err, tc.err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
