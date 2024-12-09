package common

type ApplicationState string

const (
	// Any change here should reflect in the IsValid() method too
	AppliedAppState ApplicationState = "APPLIED"

	// TODO: Remember to Reject all open applications when an Opening is closed
	RejectedAppState ApplicationState = "REJECTED"

	ShortlistedAppState ApplicationState = "SHORTLISTED"
	WithdrawnAppState   ApplicationState = "WITHDRAWN"
	ExpiredAppState     ApplicationState = "EXPIRED"
)

func (s ApplicationState) IsValid() bool {
	return s == AppliedAppState ||
		s == RejectedAppState ||
		s == ShortlistedAppState ||
		s == WithdrawnAppState ||
		s == ExpiredAppState
}
