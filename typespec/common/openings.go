package common

type OpeningState string

const (
	DraftOpening     OpeningState = "DRAFT_OPENING_STATE"
	ActiveOpening    OpeningState = "ACTIVE_OPENING_STATE"
	SuspendedOpening OpeningState = "SUSPENDED_OPENING_STATE"
	ClosedOpening    OpeningState = "CLOSED_OPENING_STATE"
)

type OpeningType string

const (
	// Any changes to here should be reflected also in the IsValid() method
	FullTimeOpening    OpeningType = "FULL_TIME_OPENING"
	PartTimeOpening    OpeningType = "PART_TIME_OPENING"
	ContractOpening    OpeningType = "CONTRACT_OPENING"
	InternshipOpening  OpeningType = "INTERNSHIP_OPENING"
	UnspecifiedOpening OpeningType = "UNSPECIFIED_OPENING"
)

func (o OpeningType) IsValid() bool {
	switch o {
	case FullTimeOpening,
		PartTimeOpening,
		ContractOpening,
		InternshipOpening,
		UnspecifiedOpening:
		return true
	}
	return false
}

type EducationLevel string

const (
	BachelorEducation    EducationLevel = "BACHELOR_EDUCATION"
	MasterEducation      EducationLevel = "MASTER_EDUCATION"
	DoctorateEducation   EducationLevel = "DOCTORATE_EDUCATION"
	NotMattersEducation  EducationLevel = "NOT_MATTERS_EDUCATION"
	UnspecifiedEducation EducationLevel = "UNSPECIFIED_EDUCATION"
)

func (e EducationLevel) IsValid() bool {
	return e == BachelorEducation || e == MasterEducation ||
		e == DoctorateEducation ||
		e == NotMattersEducation ||
		e == UnspecifiedEducation
}

type Salary struct {
	MinAmount float64  `json:"min_amount" validate:"required,min=0"`
	MaxAmount float64  `json:"max_amount" validate:"required,min=1"`
	Currency  Currency `json:"currency"   validate:"required"`
}

type FilterOpeningTagsRequest struct {
	Prefix *string `json:"prefix,omitempty"`
}

type OpeningTagID string

type OpeningTag struct {
	ID   OpeningTagID `json:"id"   validate:"required"`
	Name string       `json:"name" validate:"required"`
}

type OpeningTags struct {
	Tags []OpeningTag `json:"tags" validate:"required"`
}
