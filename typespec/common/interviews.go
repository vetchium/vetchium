package common

type InterviewState string

const (
	ScheduledInterviewState InterviewState = "SCHEDULED_INTERVIEW"
	CompletedInterviewState InterviewState = "COMPLETED_INTERVIEW"
	CancelledInterviewState InterviewState = "CANCELLED_INTERVIEW"
)

func (s InterviewState) IsValid() bool {
	switch s {
	case ScheduledInterviewState,
		CompletedInterviewState,
		CancelledInterviewState:
		return true
	}
	return false
}

type CandidacyState string

const (
	InterviewingCandidacyState CandidacyState = "INTERVIEWING"

	OfferedCandidacyState       CandidacyState = "OFFERED"
	OfferDeclinedCandidacyState CandidacyState = "OFFER_DECLINED"
	OfferAcceptedCandidacyState CandidacyState = "OFFER_ACCEPTED"

	CandidateUnsuitableCandidacyState CandidacyState = "CANDIDATE_UNSUITABLE"

	CandidateNotRespondingCandidacyState CandidacyState = "CANDIDATE_NOT_RESPONDING"
	CandidateWithdrewCandidacyState      CandidacyState = "CANDIDATE_WITHDREW"

	EmployerDefunctCandidacyState CandidacyState = "EMPLOYER_DEFUNCT"
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
	NeutralInterviewersDecision   InterviewersDecision = "NEUTRAL"
	NoInterviewersDecision        InterviewersDecision = "NO"
	StrongNoInterviewersDecision  InterviewersDecision = "STRONG_NO"
)

type RSVPStatus string

const (
	YesRSVP    RSVPStatus = "YES"
	NoRSVP     RSVPStatus = "NO"
	NotSetRSVP RSVPStatus = "NOT_SET"
)

func (s RSVPStatus) IsValidRequest() bool {
	switch s {
	case YesRSVP, NoRSVP:
		return true
	case NotSetRSVP:
		// NOT_SET is the initial value; users should not be allowed to set it
		return false
	}
	return false
}

type RSVPInterviewRequest struct {
	InterviewID string     `json:"interview_id" validate:"required"`
	RSVPStatus  RSVPStatus `json:"rsvp_status"  validate:"required,validate_rsvp_request"`
}

func (r RSVPInterviewRequest) IsValid() bool {
	switch r.RSVPStatus {
	case YesRSVP, NoRSVP, NotSetRSVP:
		return true
	}
	return false
}

type InterviewType string

const (
	InPersonInterviewType  InterviewType = "IN_PERSON"
	VideoCallInterviewType InterviewType = "VIDEO_CALL"
	TakeHomeInterviewType  InterviewType = "TAKE_HOME"
	OtherInterviewType     InterviewType = "OTHER_INTERVIEW"
)

func (i InterviewType) IsValid() bool {
	switch i {
	case InPersonInterviewType,
		VideoCallInterviewType,
		TakeHomeInterviewType,
		OtherInterviewType:
		return true
	}
	return false
}
