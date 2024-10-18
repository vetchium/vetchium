package db

import "errors"

// This file contains internal structs that can be shared between db and backend
// These are not part of the public API
// A single struct below can span across multiple db tables

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer           = errors.New("employer not found")
	ErrOrgUserAlreadyExists = errors.New("org user already exists")
	ErrNoOrgUser            = errors.New("org user not found")
)

type EmailStateChange struct {
	EmailDBKey int64
	EmailState EmailState
}

type OnboardEmailInfo struct {
	EmployerID         int64
	OnboardSecretToken string
	TokenValidMins     float64
	Email              Email
}

type OnboardInfo struct {
	EmployerID     int64
	AdminEmailAddr string
	DomainName     string
}

type OnboardReq struct {
	DomainName string
	Password   string
	Token      string
}

type OrgUserAuth struct {
	OrgUserID     int64
	OrgUserEmail  string
	EmployerID    int64
	OrgUserRole   OrgUserRole
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
