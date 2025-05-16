package libgranger

type GetEmployerCountsRequest struct {
	Domain string `json:"domain" validate:"required"`
}

type EmployerCounts struct {
	ActiveOpeningsCount    uint32 `json:"active_openings_count"`
	VerifiedEmployeesCount uint32 `json:"verified_employees_count"`
}
