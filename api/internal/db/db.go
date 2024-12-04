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

	// Used by hermione - Employer Auth related methods
	AuthOrgUser(c c.Context, sessionToken string) (OrgUserTO, error)
	CreateOrgUserToken(c.Context, EmployerTokenReq) error
	GetEmployer(c c.Context, clientID string) (Employer, error)
	GetEmployerByID(c c.Context, employerID uuid.UUID) (Employer, error)
	GetDomainNames(c c.Context, employerID uuid.UUID) ([]string, error)
	GetOrgUserAuth(c.Context, OrgUserCreds) (OrgUserAuth, error)
	GetOrgUserByTFACreds(c c.Context, tfaCode, tgt string) (OrgUserTO, error)
	InitEmployerAndDomain(c.Context, Employer, Domain) error
	InitEmployerTFA(c.Context, EmployerTFA) error
	OnboardAdmin(c.Context, OnboardReq) error
	GetHubUserByEmail(c c.Context, email string) (HubUserTO, error)
	InitHubUserTFA(c.Context, HubUserTFA) error

	// Used by hermione - Cost Center related methods
	CreateCostCenter(c.Context, v.AddCostCenterRequest) (uuid.UUID, error)
	DefunctCostCenter(c.Context, v.DefunctCostCenterRequest) error
	GetCCByName(c.Context, v.GetCostCenterRequest) (v.CostCenter, error)
	GetCostCenters(c.Context, v.GetCostCentersRequest) ([]v.CostCenter, error)
	RenameCostCenter(c.Context, v.RenameCostCenterRequest) error
	UpdateCostCenter(c.Context, v.UpdateCostCenterRequest) error

	// Used by hermione - Locations related methods
	AddLocation(c.Context, v.AddLocationRequest) (uuid.UUID, error)
	DefunctLocation(c.Context, v.DefunctLocationRequest) error
	GetLocByName(c.Context, v.GetLocationRequest) (v.Location, error)
	GetLocations(c.Context, v.GetLocationsRequest) ([]v.Location, error)
	RenameLocation(c.Context, v.RenameLocationRequest) error
	UpdateLocation(c.Context, v.UpdateLocationRequest) error

	// Used by hermione - Org users related methods
	AddOrgUser(c.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(c.Context, v.DisableOrgUserRequest) error
	EnableOrgUser(c.Context, EnableOrgUserReq) error
	FilterOrgUsers(c.Context, v.FilterOrgUsersRequest) ([]v.OrgUser, error)
	SignupOrgUser(c.Context, SignupOrgUserReq) error
	UpdateOrgUser(c.Context, v.UpdateOrgUserRequest) (uuid.UUID, error)

	// Used by granger
	CreateOnboardEmail(c.Context, OnboardEmailInfo) error
	DeQOnboard(c.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(c.Context) ([]Email, error)
	PruneTokens(c.Context) error
	UpdateEmailState(c.Context, EmailStateChange) error

	// Used by hermione - Openings related methods
	CreateOpening(c.Context, v.CreateOpeningRequest) (string, error)
	GetOpening(c.Context, v.GetOpeningRequest) (v.Opening, error)
	FilterOpenings(c.Context, v.FilterOpeningsRequest) ([]v.OpeningInfo, error)
	UpdateOpening(c.Context, v.UpdateOpeningRequest) error
	GetOpeningWatchers(
		c.Context,
		v.GetOpeningWatchersRequest,
	) ([]v.OrgUserShort, error)
	AddOpeningWatchers(c.Context, v.AddOpeningWatchersRequest) error
	RemoveOpeningWatcher(c.Context, v.RemoveOpeningWatcherRequest) error
	ChangeOpeningState(c.Context, v.ChangeOpeningStateRequest) error

	// Used by hermione - Applications related methods for employers
	GetApplicationsForEmployer(
		c.Context,
		v.GetApplicationsRequest,
	) ([]v.Application, error)
	SetApplicationColorTag(c.Context, v.SetApplicationColorTagRequest) error
	RemoveApplicationColorTag(
		c.Context,
		v.RemoveApplicationColorTagRequest,
	) error
	ShortlistApplication(c.Context, ShortlistRequest) error
	GetApplicationMailInfo(
		c c.Context,
		applicationID string,
	) (ApplicationMailInfo, error)

	// Used by hermione - for Hub users
	AuthHubUser(c c.Context, token string) (HubUserTO, error)
	ChangeHubUserPassword(c.Context, uuid.UUID, string) error
	CreateApplication(c.Context, ApplyOpeningReq) error
	CreateHubUserToken(c.Context, HubTokenReq) error
	GetHubUserByTFACreds(c.Context, string, string) (HubUserTO, error)
	FindHubOpenings(
		c.Context,
		*v.FindHubOpeningsRequest,
	) ([]v.HubOpening, error)
	GetMyHandle(c.Context) (string, error)
	InitHubUserPasswordReset(c.Context, HubUserInitPasswordReset) error
	Logout(c c.Context, token string) error
	ResetHubUserPassword(c.Context, HubUserPasswordReset) error
}
