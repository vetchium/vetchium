package vetchi

type CandidacyState string

const (
	InterviewingCandidacyState    CandidacyState = "INTERVIEWING"
	OfferedCandidacyState         CandidacyState = "OFFERED"
	OfferDeclinedCandidacyState   CandidacyState = "OFFER_DECLINED"
	OfferAcceptedCandidacyState   CandidacyState = "OFFER_ACCEPTED"
	UnsuitableCandidacyState      CandidacyState = "CANDIDATE_UNSUITABLE"
	NotRespondingCandidacyState   CandidacyState = "CANDIDATE_NOT_RESPONDING"
	EmployerDefunctCandidacyState CandidacyState = "EMPLOYER_DEFUNCT"
)
