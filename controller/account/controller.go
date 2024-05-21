package account

import (
	"net/http"
	"time"

	entity "github.com/SeongUgKim/simplebank/entity/account"
	"github.com/SeongUgKim/simplebank/lib/errorresponse"
	repository "github.com/SeongUgKim/simplebank/repository/account"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

type Controller interface {
	Create(ctx *gin.Context)
}

type controller struct {
	repository repository.Repository
}

type Params struct {
	Repository repository.Repository
}

func New(params Params) (Controller, error) {
	return &controller{
		repository: params.Repository,
	}, nil
}

func (c *controller) Create(ctx *gin.Context) {
	var req entity.CreateAccountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	uuid, err := uuid.NewV4()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
	}

	account := entity.Account{
		UUID:      uuid.String(),
		Owner:     req.Owner,
		AmountE5:  0,
		Currency:  req.Currency,
		CreatedAt: time.Now(),
	}

	insertedAccount, err := c.repository.Insert(ctx, account)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
	}

	ctx.JSON(http.StatusOK, insertedAccount)
}
