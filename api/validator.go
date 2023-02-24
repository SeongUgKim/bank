package api

import (
	"github.com/SeongUgKim/simplebank/util"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(filedLevel validator.FieldLevel) bool {
	if currency, ok := filedLevel.Field().Interface().(string); ok {
		return util.IsSupportedCurrency(currency)
	}

	return false
}
