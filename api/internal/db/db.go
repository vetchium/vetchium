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

	AddLocation(context.Context, vetchi.AddLocationRequest) (uuid.UUID, error)
	DefunctLocation(context.Context, vetchi.DefunctLocationRequest) error
	GetLocByName(
		context.Context,
		vetchi.GetLocationRequest,
	) (vetchi.Location, error)
	GetLocations(
		context.Context,
		vetchi.GetLocationsRequest,
	) ([]vetchi.Location, error)
	RenameLocation(context.Context, vetchi.RenameLocationRequest) error
	UpdateLocation(context.Context, vetchi.UpdateLocationRequest) error

	AddOrgUser(context.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(context.Context, DisableOrgUserReq) error
	EnableOrgUser(context.Context, EnableOrgUserReq) error
	FilterOrgUsers(context.Context, FilterOrgUsersReq) ([]vetchi.OrgUser, error)
	SignupOrgUser(context.Context, SignupOrgUserReq) error
	UpdateOrgUser(context.Context, UpdateOrgUserReq) (uuid.UUID, error)

	// Used by granger
	CreateOnboardEmail(context.Context, OnboardEmailInfo) error
	DeQOnboard(context.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(context.Context) ([]Email, error)
	PruneTokens(context.Context) error
	UpdateEmailState(context.Context, EmailStateChange) error

	// Opening related methods
	CreateOpening(context.Context, CreateOpeningReq) (uuid.UUID, error)
	GetOpening(context.Context, GetOpeningReq) (vetchi.Opening, error)
	FilterOpenings(context.Context, FilterOpeningsReq) ([]vetchi.Opening, error)
	UpdateOpening(context.Context, UpdateOpeningReq) error
	GetOpeningWatchers(
		context.Context,
		GetOpeningWatchersReq,
	) (vetchi.OpeningWatchers, error)
	AddOpeningWatchers(context.Context, AddOpeningWatchersReq) error
	RemoveOpeningWatcher(context.Context, RemoveOpeningWatcherReq) error
	ApproveOpeningStateChange(
		context.Context,
		ApproveOpeningStateChangeReq,
	) error
	RejectOpeningStateChange(context.Context, RejectOpeningStateChangeReq) error
}
