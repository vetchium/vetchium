package vetchi

type CountryCode string
type Currency string
type EmailAddress string
type Password string

type ValidationErrors struct {
	Errors []string `json:"errors"`
}
