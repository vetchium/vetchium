package vetchi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"

	validator "github.com/go-playground/validator/v10"
)

type Vator struct {
	validate *validator.Validate
	log      *slog.Logger
}

func InitValidator(log *slog.Logger) (*Vator, error) {
	validate := validator.New()

	err := validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if len(value) < 12 || len(value) > 64 {
				return false
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	domainReg := regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	err = validate.RegisterValidation(
		"client_id",
		func(fl validator.FieldLevel) bool {
			return domainReg.MatchString(fl.Field().String())
		},
	)
	if err != nil {
		return nil, err
	}

	// The regex taken from https://stackoverflow.com/a/67686133/153586
	emailReg := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	err = validate.RegisterValidation(
		"email",
		func(fl validator.FieldLevel) bool {
			return emailReg.MatchString(fl.Field().String())
		},
	)
	if err != nil {
		return nil, err
	}

	return &Vator{validate: validate, log: log}, nil
}

// Struct validates the struct and returns true if it is valid, otherwise it
// writes the appropriate error to the http.ResponseWriter and returns false.
// The caller should not touch the http.ResponseWriter if this returns false.
func (v *Vator) Struct(w http.ResponseWriter, i interface{}) bool {
	err := v.validate.Struct(i)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			v.log.Error("invalid validation error", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return false
		}

		var validationErrors ValidationErrors
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors.Errors = append(
				validationErrors.Errors,
				err.Tag(),
			)
		}
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(validationErrors)
		if err != nil {
			v.log.Error("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return false
		}
		return false
	}
	return true
}
