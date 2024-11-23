package vetchi

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
