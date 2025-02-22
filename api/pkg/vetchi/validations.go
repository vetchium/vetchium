package vetchi

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"reflect"
	"regexp"
	"runtime/debug"
	"time"

	validator "github.com/go-playground/validator/v10"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

type Vator struct {
	validate *validator.Validate
	log      util.Logger
}

func InitValidator(log util.Logger) (*Vator, error) {
	validate := validator.New()

	err := validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if len(value) < 12 || len(value) > 64 {
				log.Dbg("Invalid password length", "length", len(value))
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
		log.Err("failed to register domain validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"email",
		func(fl validator.FieldLevel) bool {
			value := fl.Field().String()
			if len(value) < 3 || len(value) > 254 {
				log.Dbg("Invalid email length", "length", len(value))
				return false
			}

			_, err := mail.ParseAddress(value)
			return err == nil
		},
	)
	if err != nil {
		log.Err("failed to register email validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_cc_states",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]employer.CostCenterState)
			for _, state := range states {
				switch state {
				case employer.ActiveCC:
					continue
				case employer.DefunctCC:
					continue
				default:
					log.Dbg("invalid cost center state", "state", state)
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		log.Err("failed to register cc states validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_country_code",
		func(fl validator.FieldLevel) bool {
			// TODO: Validate country code is one of the ISO 3166-1 alpha-3 codes
			countryCode := fl.Field().String()
			result := len(countryCode) == 3
			log.Dbg("country code", "code", countryCode, "result", result)
			return result
		},
	)
	if err != nil {
		log.Err("failed to register country code validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_city_aka",
		func(fl validator.FieldLevel) bool {
			cities := fl.Field().Interface().([]string)
			if len(cities) > 3 {
				log.Dbg("invalid city aka count", "count", len(cities))
				return false
			}
			for _, city := range cities {
				if len(city) < 3 || len(city) > 32 {
					log.Dbg("invalid city aka length", "city", city)
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		log.Err("failed to register city aka validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_location_state",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]employer.LocationState)
			for _, state := range states {
				switch state {
				case employer.ActiveLocation:
					continue
				case employer.DefunctLocation:
					continue
				default:
					log.Dbg("invalid location state", "state", state)
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		log.Err("failed to register location state validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_org_user_state",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]employer.OrgUserState)
			for _, state := range states {
				switch state {
				case employer.ActiveOrgUserState:
					continue
				case employer.AddedOrgUserState:
					continue
				case employer.DisabledOrgUserState:
					continue
				case employer.ReplicatedOrgUserState:
					continue
				default:
					log.Dbg("invalid org user state", "state", state)
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		log.Err("failed to register org user state validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_org_user_roles",
		func(fl validator.FieldLevel) bool {
			roles := fl.Field().Interface().(common.OrgUserRoles)
			return roles.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register org user roles validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_opening_type",
		func(fl validator.FieldLevel) bool {
			openingType, ok := fl.Field().Interface().(common.OpeningType)
			if !ok {
				log.Dbg("invalid opening type", "field", fl.Field().Interface())
				return false
			}
			return openingType.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register opening type validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_opening_states",
		func(fl validator.FieldLevel) bool {
			states := fl.Field().Interface().([]common.OpeningState)
			for _, state := range states {
				switch state {
				case common.ActiveOpening:
					continue
				case common.ClosedOpening:
					continue
				case common.DraftOpening:
					continue
				case common.SuspendedOpening:
					continue
				default:
					log.Dbg("invalid opening state", "state", state)
					return false
				}
			}
			return true
		},
	)
	if err != nil {
		log.Err("failed to register opening states validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_opening_filter_start_date",
		func(fl validator.FieldLevel) bool {
			date := fl.Field().Interface().(time.Time)
			result := !date.Before(time.Now().AddDate(-3, 0, 0))
			if !result {
				log.Dbg("invalid opening filter start date", "date", date)
			}
			return result
		},
	)
	if err != nil {
		log.Err("failed to register opening filter start date ", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_opening_filter_end_date",
		func(fl validator.FieldLevel) bool {
			date := fl.Field().Interface().(time.Time)
			return !date.Before(time.Now().AddDate(-3, 0, 0))
		},
	)
	if err != nil {
		log.Err("failed to register opening filter end date", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_domain",
		func(fl validator.FieldLevel) bool {
			domain, ok := fl.Field().Interface().(string)
			if !ok {
				log.Dbg("invalid domain", "field", fl.Field().Interface())
				return false
			}

			// Domain should be at least 4 characters (a.bc)
			if len(domain) < 4 {
				log.Dbg("domain too short", "domain", domain)
				return false
			}

			// Domain should be at most 255 characters
			if len(domain) > 255 {
				log.Dbg("domain too long", "domain", domain)
				return false
			}

			result := domainReg.MatchString(domain)

			log.Dbg(
				"Validating domain against regex",
				"domain",
				domain,
				"result",
				result,
			)
			return result
		},
	)
	if err != nil {
		log.Err("failed to register domain validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"date_after",
		func(fl validator.FieldLevel) bool {
			endDateStr := fl.Field().String()
			startDateField := fl.Parent().FieldByName(fl.Param())
			if !startDateField.IsValid() {
				log.Dbg("start date field not found", "field", fl.Param())
				return false
			}
			startDateStr := startDateField.String()

			endDate, err := time.Parse("2006-01-02", endDateStr)
			if err != nil {
				log.Dbg("failed to parse end date", "error", err)
				return false
			}

			startDate, err := time.Parse("2006-01-02", startDateStr)
			if err != nil {
				log.Dbg("failed to parse start date", "error", err)
				return false
			}

			return endDate.After(startDate) || endDate.Equal(startDate)
		},
	)
	if err != nil {
		log.Err("failed to register date_after validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_currency",
		func(fl validator.FieldLevel) bool {
			currency, ok := fl.Field().Interface().(common.Currency)
			if !ok {
				return false
			}

			// TODO: Validate currency code is one of the active ISO 4217 currency codes
			result := len(currency) == 3
			if !result {
				log.Dbg("invalid currency code", "code", currency)
			}
			return result
		},
	)
	if err != nil {
		log.Err("failed to register currency validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_timezone",
		func(fl validator.FieldLevel) bool {
			timezone, ok := fl.Field().Interface().(common.TimeZone)
			if !ok {
				return false
			}

			result := timezone.IsValid()
			if !result {
				log.Dbg("invalid timezone", "timezone", timezone)
			}
			return result
		},
	)
	if err != nil {
		log.Err("failed to register timezone validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_education_level",
		func(fl validator.FieldLevel) bool {
			educationLevel, ok := fl.Field().Interface().(common.EducationLevel)
			if !ok {
				// If the education level is not provided, it is valid
				return true
			}
			result := educationLevel.IsValid()
			if !result {
				log.Dbg("invalid education level", "level", educationLevel)
			}
			return result
		},
	)
	if err != nil {
		log.Err("failed to register education level validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_application_state",
		func(fl validator.FieldLevel) bool {
			state, ok := fl.Field().Interface().(common.ApplicationState)
			if !ok {
				return false
			}
			return state.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register application state validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_application_color_tag",
		func(fl validator.FieldLevel) bool {
			colorTag, ok := fl.Field().Interface().(employer.ApplicationColorTag)
			if !ok {
				return false
			}
			return colorTag.IsValid()
		},
	)
	if err != nil {
		log.Err("application color tag validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_candidacy_state",
		func(fl validator.FieldLevel) bool {
			state, ok := fl.Field().Interface().(common.CandidacyState)
			if !ok {
				return false
			}
			return state.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register candidacy state validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_interview_state",
		func(fl validator.FieldLevel) bool {
			state, ok := fl.Field().Interface().(common.InterviewState)
			if !ok {
				return false
			}
			return state.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register interview state validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_interview_type",
		func(fl validator.FieldLevel) bool {
			state, ok := fl.Field().Interface().(common.InterviewType)
			if !ok {
				return false
			}
			return state.IsValid()
		},
	)
	if err != nil {
		log.Err("failed to register interview type validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_rsvp_request",
		func(fl validator.FieldLevel) bool {
			rsvp, ok := fl.Field().Interface().(common.RSVPStatus)
			if !ok {
				return false
			}
			return rsvp.IsValidRequest()
		},
	)
	if err != nil {
		log.Err("failed to register rsvp request validation", "error", err)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_interviewers_decision",
		func(fl validator.FieldLevel) bool {
			decision, ok := fl.Field().Interface().(common.InterviewersDecision)
			if !ok {
				log.Dbg("invalid interviewers decision", "decision", decision)
				return false
			}

			if decision != "" {
				log.Dbg(
					"validating interviewers decision",
					"decision",
					decision,
				)
				return decision.IsValid()
			} else {
				log.Dbg("interviewers decision is empty", "decision", decision)
			}
			return true
		},
	)
	if err != nil {
		log.Err(
			"failed to register interviewers decision validation",
			"error",
			err,
		)
		return nil, err
	}

	err = validate.RegisterValidation(
		"validate_handle",
		func(fl validator.FieldLevel) bool {
			field := fl.Field().Interface()

			// First try as common.Handle
			if handle, ok := field.(common.Handle); ok {
				valid := handle.IsValid()
				log.Dbg(
					"validating as common.Handle",
					"handle",
					handle,
					"valid",
					valid,
				)
				return valid
			}

			// Then try as string
			if str, ok := field.(string); ok {
				handle := common.Handle(str)
				valid := handle.IsValid()
				log.Dbg("validating as string", "handle", str, "valid", valid)
				return valid
			}

			log.Dbg("invalid handle type", "type", reflect.TypeOf(field))
			return false
		},
	)
	if err != nil {
		log.Err("failed to register handle validation", "error", err)
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
		v.log.Err(
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
			v.log.Err("invalid validation error", "error", err)
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

		v.log.Dbg("validation failed", "failedFields", failedFields)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		err = json.NewEncoder(w).
			Encode(common.ValidationErrors{Errors: failedFields})
		if err != nil {
			v.log.Err("failed to encode validation errors", "error", err)
			// This would cause a superflous error response, but we'll log it
			http.Error(w, "", http.StatusInternalServerError)
			return false
		}
		return false
	}

	return true
}
