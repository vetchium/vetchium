package vetchi

const (
	HubBaseURL      = "https://vetchi.org"
	EmployerBaseURL = "https://employer.vetchi.org"

	ProdEnv = "prod"
	DevEnv  = "dev"
	TestEnv = "test"

	EmailFrom = "no-reply@vetchi.org"

	// TODO: Should we read these lengths from a config?
	// TODO: Should this be based on strlen instead of byte size ?
	OnBoardTokenLenBytes = 16
	SessionTokenLenBytes = 32

	SessionTokenValidMins    = 60 * 24       // 1 day
	LongTermSessionValidMins = 60 * 24 * 365 // 1 year
)
