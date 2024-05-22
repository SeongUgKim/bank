package account

import (
	"database/sql"
	"errors"
	"log"
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
	List(ctx *gin.Context)
	Fetch(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
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
		return
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
		return
	}

	ctx.JSON(http.StatusOK, insertedAccount)
}

func (c *controller) List(ctx *gin.Context) {
	accounts, err := c.repository.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}

func (c *controller) Fetch(ctx *gin.Context) {
	var req entity.GetAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println(ctx.Params)
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	if err := validateUUID(req.UUID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	account, err := c.repository.Fetch(ctx, req.UUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorresponse.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *controller) Update(ctx *gin.Context) {
	var req entity.UpdateAccountRequest
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	account, err := c.repository.Update(ctx, req.UUID, req.AmountE5)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorresponse.ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

func (c *controller) Delete(ctx *gin.Context) {
	var req entity.DeleteAccountRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	if err := validateUUID(req.UUID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorresponse.ErrorResponse(err))
		return
	}

	if err := c.repository.Delete(ctx, req.UUID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorresponse.ErrorResponse(err))
		}

		ctx.JSON(http.StatusInternalServerError, errorresponse.ErrorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, "Deleted")
}

func validateUUID(reqUUID string) error {
	if _, err := uuid.FromString(reqUUID); err != nil {
		return errors.New("invalid uuid")
	}

	return nil
}
