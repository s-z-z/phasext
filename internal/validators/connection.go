package validators

import (
	v "github.com/go-playground/validator/v10"
)

// canSsh example for custom validation.
func canSsh(fl v.FieldLevel) bool {
	return fl.Field().String() != "1.1.1.1"
}
