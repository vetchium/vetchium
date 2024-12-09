package hub

import "github.com/psankar/vetchi/typespec/common"

type AddHubCandidacyCommentRequest struct {
	CandidacyID string `json:"candidacy_id" validate:"required"`
	Comment     string `json:"comment"      validate:"required,max=2048"`
}

type MyCandidacy struct {
	CandidacyID        string                `json:"candidacy_id"`
	CompanyName        string                `json:"company_name"`
	CompanyDomain      string                `json:"company_domain"`
	OpeningID          string                `json:"opening_id"`
	OpeningTitle       string                `json:"opening_title"`
	OpeningDescription string                `json:"opening_description"`
	CandidacyState     common.CandidacyState `json:"candidacy_state"`
}
