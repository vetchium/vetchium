package vetchi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/mail"
	"reflect"
	"regexp"
	"runtime/debug"

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

	err = validate.RegisterValidation(
		"email",
		func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if len(value) < 3 || len(value) > 254 {
				return false
			}

			_, err := mail.ParseAddress(value)
			return err == nil
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_cc_states",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]CostCenterState)
			for _, state := range states {
				switch state {
				case ActiveCC:
					continue
				case DefunctCC:
					continue
				default:
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_country_code",
		func(fl validator.FieldLevel) bool {
			return len(fl.Field().String()) == 3
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_city_aka",
		func(fl validator.FieldLevel) bool {
			cities := fl.Field().Interface().([]string)
			if len(cities) > 3 {
				return false
			}
			for _, city := range cities {
				if len(city) < 3 || len(city) > 32 {
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_location_state",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]LocationState)
			for _, state := range states {
				switch state {
				case ActiveLocation:
					continue
				case DefunctLocation:
					continue
				default:
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_org_user_state",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]OrgUserState)
			for _, state := range states {
				switch state {
				case ActiveOrgUserState:
					continue
				case InvitedOrgUserState:
					continue
				case AddedOrgUserState:
					continue
				case DisabledOrgUserState:
					continue
				case ReplicatedOrgUserState:
					continue
				default:
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_org_user_roles",
		func(fl validator.FieldLevel) bool {
			roles := fl.Field().Interface().([]OrgUserRole)
			if len(roles) == 0 {
				return false
			}
			for _, role := range roles {
				switch role {
				case Admin:
					continue
				case CostCentersCRUD:
					continue
				case CostCentersViewer:
					continue
				case LocationsCRUD:
					continue
				case LocationsViewer:
					continue
				case OpeningsCRUD:
					continue
				case OpeningsViewer:
					continue
				case OrgUsersCRUD:
					continue
				case OrgUsersViewer:
					continue
				default:
					return false
				}
			}
			return true
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
	// Ensure that 'i' is a pointer to a struct
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		v.log.Error(
			"provided input is not a pointer to a struct",
			"stacktrace",
			string(debug.Stack()),
		)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return false
	}

	err := v.validate.Struct(i)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			v.log.Error("invalid validation error", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return false
		}

		// Collect all failed fields
		var failedFields []string
		structType := val.Elem().Type() // Get the actual struct type

		for _, validationErr := range err.(validator.ValidationErrors) {
			// Safely retrieve the field name via reflection
			field, found := structType.FieldByName(validationErr.StructField())
			if found {
				// Use the JSON tag name if available, otherwise fallback to the field name
				jsonTag := field.Tag.Get("json")
				if jsonTag != "" {
					failedFields = append(failedFields, jsonTag)
				} else {
					failedFields = append(failedFields, validationErr.Field())
				}
			} else {
				// Fallback to field name if reflection doesn't find the field
				failedFields = append(failedFields, validationErr.Field())
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).Encode(ValidationErrors{Errors: failedFields})
		if err != nil {
			v.log.Error("failed to encode validation errors", "error", err)
			// This would cause a superflous error response, but we'll log it
			http.Error(w, "", http.StatusInternalServerError)
			return false
		}
		return false
	}

	return true
}
