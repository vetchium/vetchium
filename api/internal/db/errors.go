package db

import (
	"errors"
)

// Ideally should be a const, but go doesn't support const errors.
var (
	ErrNoEmployer           = errors.New("employer not found")
	ErrOrgUserAlreadyExists = errors.New("org user already exists")
	ErrNoOrgUser            = errors.New("org user not found")
	ErrDupCostCenterName    = errors.New("duplicate cost center name")
	ErrNoCostCenter         = errors.New("cost center not found")
)
