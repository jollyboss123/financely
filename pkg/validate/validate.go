package validate

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func Validate(v *validator.Validate, generic any) []string {
	err := v.Struct(generic)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return nil
		}

		var errs []string
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Sprintf("%s is %s with type %s", err.StructField(), err.Tag(), err.Type()))
		}

		return errs
	}
	return nil
}
