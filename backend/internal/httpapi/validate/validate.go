package validate

import "github.com/go-playground/validator/v10"

var v = validator.New()

func Struct(s any) error {
	return v.Struct(s)
}
