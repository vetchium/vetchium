package vetchi

import (
	validator "github.com/go-playground/validator/v10"
)

func InitValidator() (*validator.Validate, error) {
	validate := validator.New()

	err := validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if len(value) < 8 || len(value) > 32 {
				return false
			}

			return true
		},
	)
	if err != nil {
		return nil, err
	}

	return validate, nil
}
