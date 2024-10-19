package db

import (
	"context"

	"github.com/google/uuid"
)

// Create a new qo.Struct for passing parameters to the db functions, if any
// function declaration goes more than 80 characters or multiple lines.

// Do not name parameters when passing objects. Name parameters when passing
// primitive data types.

type DB interface {
	// Used by hermione and granger

	// Used by hermione
	InitEmployerTFA(context.Context, EmployerTFA) error
	InitEmployerAndDomain(context.Context, Employer, Domain) error
	GetEmployer(c context.Context, clientID string) (Employer, error)
	GetOrgUserAuth(context.Context, OrgUserCreds) (OrgUserAuth, error)
	GetOrgUserByToken(c context.Context, tfaCode, tgt string) (OrgUser, error)
	OnboardAdmin(context.Context, OnboardReq) error
	CreateCostCenter(context.Context, CostCenterReq) (uuid.UUID, error)
	CreateOrgUserToken(context.Context, OrgUserToken) error

	// Used by granger
	CreateOnboardEmail(context.Context, OnboardEmailInfo) error
	GetOldestUnsentEmails(context.Context) ([]Email, error)
	PruneTokens(context.Context) error
	UpdateEmailState(context.Context, EmailStateChange) error
	DeQOnboard(context.Context) (*OnboardInfo, error)
}
