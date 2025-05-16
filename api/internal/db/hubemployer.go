package db

type EmployerDetailsForHub struct {
	// TODO: We will add more fields as needed in future
	EmployerID    string
	Name          string
	PrimaryDomain string
	OtherDomains  []string
}
