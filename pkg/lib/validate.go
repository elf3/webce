package lib

import "github.com/go-playground/validator/v10"

// Validate the fields.
func Validate(u interface{}) error {
	validate := validator.New()
	return validate.Struct(u)
}
