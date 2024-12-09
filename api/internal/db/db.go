package db

import (
	c "context"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
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
	CreateCostCenter(
		c.Context,
		employer.AddCostCenterRequest,
	) (uuid.UUID, error)
	DefunctCostCenter(c.Context, employer.DefunctCostCenterRequest) error
	GetCCByName(
		c.Context,
		employer.GetCostCenterRequest,
	) (employer.CostCenter, error)
	GetCostCenters(
		c.Context,
		employer.GetCostCentersRequest,
	) ([]employer.CostCenter, error)
	RenameCostCenter(c.Context, employer.RenameCostCenterRequest) error
	UpdateCostCenter(c.Context, employer.UpdateCostCenterRequest) error

	// Used by hermione - Locations related methods
	AddLocation(c.Context, employer.AddLocationRequest) (uuid.UUID, error)
	DefunctLocation(c.Context, employer.DefunctLocationRequest) error
	GetLocByName(
		c.Context,
		employer.GetLocationRequest,
	) (employer.Location, error)
	GetLocations(
		c.Context,
		employer.GetLocationsRequest,
	) ([]employer.Location, error)
	RenameLocation(c.Context, employer.RenameLocationRequest) error
	UpdateLocation(c.Context, employer.UpdateLocationRequest) error

	// Used by hermione - Org users related methods
	AddOrgUser(c.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(c.Context, employer.DisableOrgUserRequest) error
	EnableOrgUser(c.Context, EnableOrgUserReq) error
	FilterOrgUsers(
		c.Context,
		employer.FilterOrgUsersRequest,
	) ([]employer.OrgUser, error)
	SignupOrgUser(c.Context, SignupOrgUserReq) error
	UpdateOrgUser(c.Context, employer.UpdateOrgUserRequest) (uuid.UUID, error)

	// Used by granger
	CreateOnboardEmail(c.Context, OnboardEmailInfo) error
	DeQOnboard(c.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(c.Context) ([]Email, error)
	PruneTokens(c.Context) error
	UpdateEmailState(c.Context, EmailStateChange) error

	// Used by hermione - Openings related methods
	CreateOpening(c.Context, employer.CreateOpeningRequest) (string, error)
	GetOpening(c.Context, employer.GetOpeningRequest) (employer.Opening, error)
	FilterOpenings(
		c.Context,
		employer.FilterOpeningsRequest,
	) ([]employer.OpeningInfo, error)
	UpdateOpening(c.Context, employer.UpdateOpeningRequest) error
	GetOpeningWatchers(
		c.Context,
		employer.GetOpeningWatchersRequest,
	) ([]employer.OrgUserShort, error)
	AddOpeningWatchers(c.Context, employer.AddOpeningWatchersRequest) error
	RemoveOpeningWatcher(c.Context, employer.RemoveOpeningWatcherRequest) error
	ChangeOpeningState(c.Context, employer.ChangeOpeningStateRequest) error

	// Used by hermione - Applications related methods for employers
	GetApplicationsForEmployer(
		c.Context,
		employer.GetApplicationsRequest,
	) ([]employer.Application, error)
	SetApplicationColorTag(
		c.Context,
		employer.SetApplicationColorTagRequest,
	) error
	RemoveApplicationColorTag(
		c.Context,
		employer.RemoveApplicationColorTagRequest,
	) error
	ShortlistApplication(c.Context, ShortlistRequest) error
	RejectApplication(c.Context, RejectApplicationRequest) error
	GetApplicationMailInfo(
		c c.Context,
		applicationID string,
	) (ApplicationMailInfo, error)

	// Used by hermione - for Hub users
	AuthHubUser(c c.Context, token string) (HubUserTO, error)
	ChangeHubUserPassword(c.Context, uuid.UUID, string) error
	CreateApplication(c.Context, ApplyOpeningReq) error
	MyApplications(
		c.Context,
		hub.MyApplicationsRequest,
	) ([]hub.HubApplication, error)
	CreateHubUserToken(c.Context, HubTokenReq) error
	GetHubUserByTFACreds(c.Context, string, string) (HubUserTO, error)
	FindHubOpenings(
		c.Context,
		*hub.FindHubOpeningsRequest,
	) ([]hub.HubOpening, error)
	GetMyHandle(c.Context) (string, error)
	InitHubUserPasswordReset(c.Context, HubUserInitPasswordReset) error
	Logout(c c.Context, token string) error
	ResetHubUserPassword(c.Context, HubUserPasswordReset) error
}
