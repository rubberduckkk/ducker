package rest

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Validator interface {
	IsValid() bool
}

func registerValidation() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("expect *validator.Validate")
	}

	if err := v.RegisterValidation("query_order", customerValidator); err != nil {
		panic(fmt.Errorf("register query order validation failed. %w", err))
	}
}

func customerValidator(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(Validator)
	return value.IsValid()
}
