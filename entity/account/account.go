package account

import "time"

type Account struct {
	UUID      string    `db:"uuid"`
	Owner     string    `db:"owner"`
	AmountE5  int64     `db:"amount_e5"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}
