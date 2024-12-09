package common

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
