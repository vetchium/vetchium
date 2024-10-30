package db

import "github.com/google/uuid"

type CCenterReq struct {
	Name       string
	Notes      string
	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type GetCCByNameReq struct {
	Name       string
	EmployerID uuid.UUID
}

type RenameCCReq struct {
	OldName    string
	NewName    string
	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type UpdateCCReq struct {
	Name       string
	Notes      string
	EmployerID uuid.UUID

	// Currently unused, but will be used in the future for audit logs
	OrgUserID uuid.UUID
}

type CCentersList struct {
	EmployerID uuid.UUID
	States     []string

	Limit         int
	PaginationKey string
}

type DefunctCCReq struct {
	EmployerID uuid.UUID
	Name       string
}
