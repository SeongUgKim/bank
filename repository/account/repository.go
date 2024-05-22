package account

import (
	"context"
	"database/sql"
	"time"

	db "github.com/SeongUgKim/simplebank/db/postgres"
	entity "github.com/SeongUgKim/simplebank/entity/account"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	isnertQuery = `
	INSERT INTO accounts (
		uuid,
		owner,
		amount_e5, 
		currency,
		created_at
	) VALUES (
		:uuid,
		:owner,
		:amount_e5,
		:currency,
		:created_at
	)
	`
	updateQuery = "UPDATE accounts SET amount_e5 = $1, created_at = $2 WHERE uuid = $3"
	fetchQuery  = `SELECT * FROM accounts WHERE uuid = $1`
	listQuery   = `SELECT * FROM accounts ORDER BY created_at`
	deleteQuery = `DELETE FROM accounts WHERE uuid = $1`
)

type Repository interface {
	Insert(ctx context.Context, account entity.Account) (entity.Account, error)
	Fetch(ctx context.Context, accountUUID string) (entity.Account, error)
	List(ctx context.Context) ([]entity.Account, error)
	Update(ctx context.Context, accountUUID string, amountE5 int64) (entity.Account, error)
	Delete(ctx context.Context, accountUUID string) error
}

type repository struct {
	db *db.PostgresDB
}

type Params struct {
	DB *db.PostgresDB
}

func New(params Params) (Repository, error) {
	return &repository{
		db: params.DB,
	}, nil
}

func (r *repository) Insert(ctx context.Context, account entity.Account) (entity.Account, error) {
	tx, err := r.db.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return entity.Account{}, errors.Wrapf(err, "failed to begin transaction")
	}

	if err := insertInTx(ctx, tx, account); err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	insertedAccount, err := readInTx(ctx, tx, account.UUID)
	if err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	return insertedAccount, nil
}

func (r *repository) Fetch(ctx context.Context, accountUUID string) (entity.Account, error) {
	var account entity.Account
	if err := r.db.DB.GetContext(ctx, &account, fetchQuery, accountUUID); err != nil {
		return entity.Account{}, errors.Wrapf(err, "failed to fetch account from db")
	}

	return account, nil
}

func (r *repository) List(ctx context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	if err := r.db.DB.SelectContext(ctx, &accounts, listQuery); err != nil {
		return nil, errors.Wrapf(err, "failed to fetch accounts from db")
	}

	return accounts, nil
}

func (r *repository) Update(ctx context.Context, accountUUID string, amountE5 int64) (entity.Account, error) {
	tx, err := r.db.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return entity.Account{}, errors.Wrapf(err, "failed to begin transaction")
	}

	account, err := readInTx(ctx, tx, accountUUID)
	if err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	newAccount := entity.Account{
		UUID:      accountUUID,
		Owner:     account.Owner,
		AmountE5:  amountE5,
		Currency:  account.Currency,
		CreatedAt: time.Now(),
	}

	if err := updateInTx(ctx, tx, newAccount); err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	updatedAccount, err := readInTx(ctx, tx, account.UUID)
	if err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return entity.Account{}, db.Rollback(tx, err)
	}

	return updatedAccount, nil

}

func (r *repository) Delete(ctx context.Context, accountUUID string) error {
	tx, err := r.db.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		return errors.Wrapf(err, "failed to begin transaction")
	}

	if _, err := readInTx(ctx, tx, accountUUID); err != nil {
		return db.Rollback(tx, err)
	}

	if err := deleteInTx(ctx, tx, accountUUID); err != nil {
		return db.Rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return db.Rollback(tx, err)
	}

	return nil
}

func updateInTx(ctx context.Context, tx *sqlx.Tx, account entity.Account) error {
	qry, args, err := sqlx.In(updateQuery, account.AmountE5, account.CreatedAt, account.UUID)
	if err != nil {
		return errors.Wrapf(err, "failed to create update query")
	}

	if _, err := tx.ExecContext(ctx, qry, args...); err != nil {
		return errors.Wrapf(err, "failed to update account in db")
	}

	return nil
}

func deleteInTx(ctx context.Context, tx *sqlx.Tx, accountUUID string) error {
	qry, args, err := sqlx.In(deleteQuery, accountUUID)
	if err != nil {
		return errors.Wrapf(err, "failed to create delete query")
	}

	if _, err := tx.ExecContext(ctx, qry, args...); err != nil {
		return errors.Wrapf(err, "failed to delete account in db")
	}

	return nil
}

func insertInTx(ctx context.Context, tx *sqlx.Tx, account entity.Account) error {
	if _, err := tx.NamedExecContext(ctx, isnertQuery, account); err != nil {
		return errors.Wrapf(err, "failed to insert account in db")
	}

	return nil
}

func readInTx(ctx context.Context, tx *sqlx.Tx, accountUUID string) (entity.Account, error) {
	var account entity.Account
	if err := tx.GetContext(ctx, &account, fetchQuery, accountUUID); err != nil {
		return entity.Account{}, errors.Wrapf(err, "failed to fetch account from db")
	}

	return account, nil
}
