package common

type InterviewState string

const (
	ScheduledInterviewState         InterviewState = "SCHEDULED"
	CompletedInterviewState         InterviewState = "COMPLETED"
	CandidateWithdrewInterviewState InterviewState = "CANDIDATE_WITHDREW"
	EmployerWithdrewInterviewState  InterviewState = "EMPLOYER_WITHDREW"
)

func (s InterviewState) IsValid() bool {
	switch s {
	case ScheduledInterviewState,
		CompletedInterviewState,
		CandidateWithdrewInterviewState,
		EmployerWithdrewInterviewState:
		return true
	}
	return false
}

type CandidacyState string

const (
	InterviewingCandidacyState           CandidacyState = "INTERVIEWING"
	OfferedCandidacyState                CandidacyState = "OFFERED"
	OfferDeclinedCandidacyState          CandidacyState = "OFFER_DECLINED"
	OfferAcceptedCandidacyState          CandidacyState = "OFFER_ACCEPTED"
	CandidateUnsuitableCandidacyState    CandidacyState = "CANDIDATE_UNSUITABLE"
	CandidateNotRespondingCandidacyState CandidacyState = "CANDIDATE_NOT_RESPONDING"
	EmployerDefunctCandidacyState        CandidacyState = "EMPLOYER_DEFUNCT"
)

func (s CandidacyState) IsValid() bool {
	switch s {
	case InterviewingCandidacyState,
		OfferedCandidacyState,
		OfferDeclinedCandidacyState,
		OfferAcceptedCandidacyState,
		CandidateUnsuitableCandidacyState,
		CandidateNotRespondingCandidacyState,
		EmployerDefunctCandidacyState:
		return true
	}
	return false
}

type InterviewersDecision string

const (
	StrongYesInterviewersDecision InterviewersDecision = "STRONG_YES"
	YesInterviewersDecision       InterviewersDecision = "YES"
	NoInterviewersDecision        InterviewersDecision = "NO"
	StrongNoInterviewersDecision  InterviewersDecision = "STRONG_NO"
)

type RSVPStatus string

const (
	YesRSVP    RSVPStatus = "YES"
	NoRSVP     RSVPStatus = "NO"
	NotSetRSVP RSVPStatus = "NOT_SET"
)

func (s RSVPStatus) IsValid() bool {
	switch s {
	case YesRSVP, NoRSVP, NotSetRSVP:
		return true
	}
	return false
}

type RSVPInterviewRequest struct {
	InterviewID string     `json:"interview_id" validate:"required"`
	RSVPStatus  RSVPStatus `json:"rsvp_status"  validate:"required,validate_rsvp_status"`
}

func (r RSVPInterviewRequest) IsValid() bool {
	switch r.RSVPStatus {
	case YesRSVP, NoRSVP, NotSetRSVP:
		return true
	}
	return false
}
