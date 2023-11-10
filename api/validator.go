package api

import (
	"github.com/AbdelilahOu/GoThingy/utils"
	"github.com/go-playground/validator/v10"
)

var validCurrency validator.Func = func(feildLevel validator.FieldLevel) bool {
	if currency, ok := feildLevel.Field().Interface().(string); ok {
		return utils.IsSupportedCurrency(currency)
	}
	return false
}
