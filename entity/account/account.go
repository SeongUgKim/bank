package account

import (
	"time"

	_ "github.com/gin-gonic/gin"
)

type Account struct {
	UUID      string    `db:"uuid"`
	Owner     string    `db:"owner"`
	AmountE5  int64     `db:"amount_e5"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

type CreateAccountRequest struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,oneof=USD EUR"`
}

type GetAccountRequest struct {
	UUID string `uri:"uuid" binding:"required"`
}

type UpdateAccountRequest struct {
	UUID     string `json:"uuid" binding:"required"`
	AmountE5 int64  `json:"amount_e5" binding:"required"`
}

type DeleteAccountRequest struct {
	UUID string `uri:"uuid" binding:"required"`
}
