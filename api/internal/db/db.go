package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

// Create a new qo.Struct for passing parameters to the db functions, if any
// function declaration goes more than 80 characters or multiple lines.

// Do not name parameters when passing objects. Name parameters when passing
// primitive data types.

type DB interface {
	// Used by hermione and granger

	// Used by hermione
	AuthOrgUser(c context.Context, sessionToken string) (OrgUserTO, error)
	CreateOrgUserToken(context.Context, TokenReq) error
	GetEmployer(c context.Context, clientID string) (Employer, error)
	GetEmployerByID(c context.Context, employerID uuid.UUID) (Employer, error)
	GetDomainNames(c context.Context, employerID uuid.UUID) ([]string, error)
	GetOrgUserAuth(context.Context, OrgUserCreds) (OrgUserAuth, error)
	GetOrgUserByToken(c context.Context, tfaCode, tgt string) (OrgUserTO, error)
	InitEmployerAndDomain(context.Context, Employer, Domain) error
	InitEmployerTFA(context.Context, EmployerTFA) error
	OnboardAdmin(context.Context, OnboardReq) error

	CreateCostCenter(context.Context, CCenterReq) (uuid.UUID, error)
	DefunctCostCenter(context.Context, DefunctCCReq) error
	GetCCByName(context.Context, GetCCByNameReq) (vetchi.CostCenter, error)
	GetCostCenters(context.Context, CCentersList) ([]vetchi.CostCenter, error)
	RenameCostCenter(context.Context, RenameCCReq) error
	UpdateCostCenter(context.Context, UpdateCCReq) error

	AddLocation(context.Context, AddLocationReq) (uuid.UUID, error)
	DefunctLocation(context.Context, DefunctLocationReq) error
	GetLocByName(context.Context, GetLocByNameReq) (vetchi.Location, error)
	GetLocations(context.Context, GetLocationsReq) ([]vetchi.Location, error)
	RenameLocation(context.Context, RenameLocationReq) error
	UpdateLocation(context.Context, UpdateLocationReq) error

	AddOrgUser(context.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(context.Context, DisableOrgUserReq) error
	EnableOrgUser(context.Context, EnableOrgUserReq) error
	FilterOrgUsers(context.Context, FilterOrgUsersReq) ([]vetchi.OrgUser, error)
	UpdateOrgUser(context.Context, UpdateOrgUserReq) (uuid.UUID, error)

	// Used by granger
	CreateOnboardEmail(context.Context, OnboardEmailInfo) error
	DeQOnboard(context.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(context.Context) ([]Email, error)
	PruneTokens(context.Context) error
	UpdateEmailState(context.Context, EmailStateChange) error
}
