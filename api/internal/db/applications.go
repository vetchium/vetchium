package db

// ResumeDetails contains the information needed to retrieve a resume file
type ResumeDetails struct {
	SHA           string
	HubUserHandle string
	ApplicationID string
}

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
