package account

import (
	controller "github.com/SeongUgKim/simplebank/controller/account"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	Start(address string) error
}

type handler struct {
	controller controller.Controller
	router     *gin.Engine
}

type Params struct {
	Controller controller.Controller
}

func New(params Params) (Handler, error) {
	handler := handler{
		controller: params.Controller,
	}

	router := gin.Default()
	router.POST("/accounts", handler.controller.Create)
	router.GET("/accounts", handler.controller.List)
	router.GET("/accounts/:uuid", handler.controller.Fetch)
	router.PATCH("/accounts", handler.controller.Update)
	router.DELETE("/accounts/:uuid", handler.controller.Delete)
	handler.router = router

	return &handler, nil
}

func (h *handler) Start(address string) error {
	return h.router.Run(address)
}
