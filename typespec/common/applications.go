package common

type ApplicationState string

const (
	AppliedAppState     ApplicationState = "APPLIED"
	RejectedAppState    ApplicationState = "REJECTED"
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

type ApplicationColorTag string

const (
	GreenApplicationColorTag  ApplicationColorTag = "GREEN"
	YellowApplicationColorTag ApplicationColorTag = "YELLOW"
	RedApplicationColorTag    ApplicationColorTag = "RED"
)

func (c ApplicationColorTag) IsValid() bool {
	return c == GreenApplicationColorTag ||
		c == YellowApplicationColorTag ||
		c == RedApplicationColorTag
}
