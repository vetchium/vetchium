package db

type ShortlistRequest struct {
	ApplicationID string
	OpeningID     string
	CandidacyID   string
	Email         Email
}

type RejectApplicationRequest struct {
	ApplicationID string
	Email         Email
}
