package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PostgresDB struct {
	DB *sqlx.DB
}

func New(driver, path string) (*PostgresDB, error) {
	db, err := sqlx.Connect(driver, path)
	if err != nil {
		return nil, err
	}

	return &PostgresDB{DB: db}, nil
}

func Rollback(tx *sqlx.Tx, err error) error {
	if rollbackErr := tx.Rollback(); rollbackErr != nil {
		return errors.Wrapf(
			errors.Wrapf(err, rollbackErr.Error()),
			"failed to rollback transaction after error",
		)
	}

	return err
}
