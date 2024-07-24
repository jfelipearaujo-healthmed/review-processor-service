package validator

import (
	goValidator "github.com/go-playground/validator/v10"
	"github.com/klassmann/cpfcnpj"
)

var val *goValidator.Validate

func init() {
	val = goValidator.New(goValidator.WithRequiredStructEnabled())
	if err := val.RegisterValidation("cpfcnpj", validateCpfCnpj); err != nil {
		panic(err)
	}
}

func Validate(i interface{}) error {
	return val.Struct(i)
}

func validateCpfCnpj(f goValidator.FieldLevel) bool {
	ok := cpfcnpj.ValidateCPF(f.Field().String())

	if !ok {
		ok = cpfcnpj.ValidateCNPJ(f.Field().String())
	}

	return ok
}
