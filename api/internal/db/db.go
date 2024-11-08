package db

import (
	c "context"

	"github.com/google/uuid"
	v "github.com/psankar/vetchi/api/pkg/vetchi"
)

// Do not name parameters when passing objects. Name parameters when passing
// primitive data types. Try to keep within 80 characters.

type DB interface {
	// Used by hermione and granger

	// Used by hermione
	AuthOrgUser(c c.Context, sessionToken string) (OrgUserTO, error)
	CreateOrgUserToken(c.Context, TokenReq) error
	GetEmployer(c c.Context, clientID string) (Employer, error)
	GetEmployerByID(c c.Context, employerID uuid.UUID) (Employer, error)
	GetDomainNames(c c.Context, employerID uuid.UUID) ([]string, error)
	GetOrgUserAuth(c.Context, OrgUserCreds) (OrgUserAuth, error)
	GetOrgUserByToken(c c.Context, tfaCode, tgt string) (OrgUserTO, error)
	InitEmployerAndDomain(c.Context, Employer, Domain) error
	InitEmployerTFA(c.Context, EmployerTFA) error
	OnboardAdmin(c.Context, OnboardReq) error

	CreateCostCenter(c.Context, v.AddCostCenterRequest) (uuid.UUID, error)
	DefunctCostCenter(c.Context, v.DefunctCostCenterRequest) error
	GetCCByName(c.Context, v.GetCostCenterRequest) (v.CostCenter, error)
	GetCostCenters(c.Context, v.GetCostCentersRequest) ([]v.CostCenter, error)
	RenameCostCenter(c.Context, v.RenameCostCenterRequest) error
	UpdateCostCenter(c.Context, v.UpdateCostCenterRequest) error

	// Locations related methods
	AddLocation(c.Context, v.AddLocationRequest) (uuid.UUID, error)
	DefunctLocation(c.Context, v.DefunctLocationRequest) error
	GetLocByName(c.Context, v.GetLocationRequest) (v.Location, error)
	GetLocations(c.Context, v.GetLocationsRequest) ([]v.Location, error)
	RenameLocation(c.Context, v.RenameLocationRequest) error
	UpdateLocation(c.Context, v.UpdateLocationRequest) error

	// Org users related methods
	AddOrgUser(c.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(c.Context, DisableOrgUserReq) error
	EnableOrgUser(c.Context, EnableOrgUserReq) error
	FilterOrgUsers(c.Context, FilterOrgUsersReq) ([]v.OrgUser, error)
	SignupOrgUser(c.Context, SignupOrgUserReq) error
	UpdateOrgUser(c.Context, UpdateOrgUserReq) (uuid.UUID, error)

	// Used by granger
	CreateOnboardEmail(c.Context, OnboardEmailInfo) error
	DeQOnboard(c.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(c.Context) ([]Email, error)
	PruneTokens(c.Context) error
	UpdateEmailState(c.Context, EmailStateChange) error

	// Openings related methods
	CreateOpening(c.Context, v.CreateOpeningRequest) (uuid.UUID, error)
	GetOpening(c.Context, v.GetOpeningRequest) (v.Opening, error)
	FilterOpenings(c.Context, v.FilterOpeningsRequest) ([]v.Opening, error)
	UpdateOpening(c.Context, v.UpdateOpeningRequest) error
	GetOpeningWatchers(
		c.Context,
		v.GetOpeningWatchersRequest,
	) (v.OpeningWatchers, error)
	AddOpeningWatchers(c.Context, v.AddOpeningWatchersRequest) error
	RemoveOpeningWatcher(c.Context, v.RemoveOpeningWatcherRequest) error
	ApproveOpeningStateChange(
		c.Context,
		v.ApproveOpeningStateChangeRequest,
	) error
	RejectOpeningStateChange(c.Context, v.RejectOpeningStateChangeRequest) error
}
