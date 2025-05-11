package validators

import v "github.com/go-playground/validator/v10"

func RegisteValidator(v *v.Validate) {
	v.RegisterValidation("canSsh", canSsh)
}
