package db

import "errors"

// This file contains the struct defintions for data that should be retrieved
// from the database and passed on to the backend.

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer           = errors.New("employer not found")
	ErrOrgUserAlreadyExists = errors.New("org user already exists")
	ErrNoOrgUser            = errors.New("org user not found")
)

type OnboardInfo struct {
	EmployerID     int64
	AdminEmailAddr string
	DomainName     string
}

type OrgUserAuth struct {
	OrgUserID     int64
	EmployerID    int64
	OrgUserRole   OrgUserRole
	PasswordHash  string
	EmployerState EmployerState
	OrgUserState  OrgUserState
}
