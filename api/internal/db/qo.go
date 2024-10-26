package db

import (
	"errors"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// This file contains internal structs that can be shared between db and backend
// These are not part of the public API
// A single struct below can span across multiple db tables

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer           = errors.New("employer not found")
	ErrOrgUserAlreadyExists = errors.New("org user already exists")
	ErrNoOrgUser            = errors.New("org user not found")
	ErrDupCostCenterName    = errors.New("duplicate cost center name")
	ErrNoCostCenter         = errors.New("cost center not found")
)

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

type DefunctReq struct {
	EmployerID uuid.UUID
	Name       string
}

type EmailStateChange struct {
	EmailDBKey uuid.UUID
	EmailState EmailState
}

type OnboardEmailInfo struct {
	EmployerID         uuid.UUID
	OnboardSecretToken string
	TokenValidMins     float64
	Email              Email
}

type OnboardInfo struct {
	EmployerID     uuid.UUID
	AdminEmailAddr string
	DomainName     string
}

type OnboardReq struct {
	DomainName string
	Password   string
	Token      string
}

type OrgUserAuth struct {
	OrgUserID     uuid.UUID
	OrgUserEmail  string
	EmployerID    uuid.UUID
	OrgUserRoles  []vetchi.OrgUserRole
	PasswordHash  string
	EmployerState EmployerState
	OrgUserState  OrgUserState
}

type OrgUserCreds struct {
	ClientID string
	Email    string
}

type EmployerTFA struct {
	EmailToken OrgUserToken
	TGToken    OrgUserToken
	Email      Email
}
