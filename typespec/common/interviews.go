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
