package db

import (
	"context"
	"time"

	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
	"github.com/vetchium/vetchium/typespec/hub"

	"github.com/google/uuid"
)

// Do not name parameters when passing objects. Name parameters when passing
// primitive data types. Try to keep within 80 characters.

type OfferToCandidateReq struct {
	CandidacyID string
	Comment     string
	Email       Email
}

type CandidateInfo struct {
	CandidateName  string
	CandidateEmail string
	CompanyName    string
	OpeningTitle   string
}

type UpdateOfficialEmailVerificationCodeReq struct {
	Email   string
	Code    string
	HubUser HubUserTO
}

type OfficialEmail struct {
	Email            string
	LastVerifiedAt   *time.Time
	VerifyInProgress bool
}

type StaleFile struct {
	ID       uuid.UUID
	FilePath string
}

type HubUserContact struct {
	Handle   string
	FullName string
	Email    string
}

// ApplicationScoringModel represents an AI model used for resume scoring
type ApplicationScoringModel struct {
	ModelName   string
	Description string
	IsActive    bool
	CreatedAt   time.Time
}

// ApplicationScore represents a score given to an application by an AI model
type ApplicationScore struct {
	ID            string
	ApplicationID string
	ModelName     string
	Score         int
	CreatedAt     time.Time
}

// Application represents an application with minimal fields needed for resume scoring
type ApplicationForScoring struct {
	ApplicationID string
	ResumeSHA     string
}

// UnscoredApplicationBatch represents a batch of unscored applications for a single opening
type UnscoredApplicationBatch struct {
	EmployerID   string
	OpeningID    string
	JD           string
	Applications []ApplicationForScoring // max 10 elements
}

// TODO: We should group this interface better
type DB interface {
	// Used by hermione and granger

	// Used by hermione - Employer Auth related methods
	AuthOrgUser(c context.Context, sessionToken string) (OrgUserTO, error)
	CreateOrgUserToken(context.Context, EmployerTokenReq) error
	GetEmployer(c context.Context, clientID string) (Employer, error)
	GetEmployerByID(c context.Context, employerID uuid.UUID) (Employer, error)
	GetDomainNames(c context.Context, employerID uuid.UUID) ([]string, error)
	GetOrgUserAuth(context.Context, OrgUserCreds) (OrgUserAuth, error)
	GetOrgUserByTFACreds(
		c context.Context,
		tfaCode, tgt string,
	) (OrgUserTO, error)
	InitEmployerAndDomain(context.Context, Employer, Domain) error
	InitEmployerTFA(context.Context, EmployerTFA) error
	OnboardAdmin(context.Context, OnboardReq) error
	GetHubUserByEmail(c context.Context, email string) (HubUserTO, error)
	InitHubUserTFA(context.Context, HubUserTFA) error

	// Used by hermione - Cost Center related methods
	CreateCostCenter(
		context.Context,
		employer.AddCostCenterRequest,
	) (uuid.UUID, error)
	DefunctCostCenter(context.Context, employer.DefunctCostCenterRequest) error
	GetCCByName(
		context.Context,
		employer.GetCostCenterRequest,
	) (employer.CostCenter, error)
	GetCostCenters(
		context.Context,
		employer.GetCostCentersRequest,
	) ([]employer.CostCenter, error)
	RenameCostCenter(context.Context, employer.RenameCostCenterRequest) error
	UpdateCostCenter(context.Context, employer.UpdateCostCenterRequest) error

	// Used by hermione - Locations related methods
	AddLocation(context.Context, employer.AddLocationRequest) (uuid.UUID, error)
	DefunctLocation(context.Context, employer.DefunctLocationRequest) error
	GetLocByName(
		context.Context,
		employer.GetLocationRequest,
	) (employer.Location, error)
	GetLocations(
		context.Context,
		employer.GetLocationsRequest,
	) ([]employer.Location, error)
	RenameLocation(context.Context, employer.RenameLocationRequest) error
	UpdateLocation(context.Context, employer.UpdateLocationRequest) error

	// Used by hermione - Org users related methods
	AddOrgUser(context.Context, AddOrgUserReq) (uuid.UUID, error)
	DisableOrgUser(context.Context, employer.DisableOrgUserRequest) error
	EnableOrgUser(context.Context, EnableOrgUserReq) error
	FilterOrgUsers(
		context.Context,
		employer.FilterOrgUsersRequest,
	) ([]employer.OrgUser, error)
	SignupOrgUser(context.Context, SignupOrgUserReq) error
	UpdateOrgUser(
		context.Context,
		employer.UpdateOrgUserRequest,
	) (uuid.UUID, error)
	GetOrgUserByEmail(context.Context, string) (OrgUserTO, error)

	// Used by granger
	CreateOnboardEmail(context.Context, OnboardEmailInfo) error
	DeQOnboard(context.Context) (*OnboardInfo, error)
	GetOldestUnsentEmails(context.Context) ([]Email, error)
	PruneTokens(context.Context) error
	UpdateEmailState(context.Context, EmailStateChange) error

	// Used by hermione - Openings related methods
	CreateOpening(
		context.Context,
		employer.CreateOpeningRequest,
	) (string, error)
	GetOpening(
		context.Context,
		employer.GetOpeningRequest,
	) (employer.Opening, error)
	FilterOpenings(
		context.Context,
		employer.FilterOpeningsRequest,
	) ([]employer.OpeningInfo, error)
	UpdateOpening(context.Context, employer.UpdateOpeningRequest) error
	GetOpeningWatchers(
		context.Context,
		employer.GetOpeningWatchersRequest,
	) ([]employer.OrgUserShort, error)
	AddOpeningWatchers(
		context.Context,
		employer.AddOpeningWatchersRequest,
	) error
	RemoveOpeningWatcher(
		context.Context,
		employer.RemoveOpeningWatcherRequest,
	) error
	ChangeOpeningState(
		context.Context,
		employer.ChangeOpeningStateRequest,
	) error

	// Used by hermione - Applications related methods for employers
	GetApplicationsForEmployer(
		context.Context,
		employer.GetApplicationsRequest,
	) ([]employer.Application, error)
	GetResumeDetails(
		context.Context,
		employer.GetResumeRequest,
	) (ResumeDetails, error)
	SetApplicationColorTag(
		context.Context,
		employer.SetApplicationColorTagRequest,
	) error
	RemoveApplicationColorTag(
		context.Context,
		employer.RemoveApplicationColorTagRequest,
	) error
	ShortlistApplication(context.Context, ShortlistRequest) error
	RejectApplication(context.Context, RejectApplicationRequest) error
	GetApplicationMailInfo(
		c context.Context,
		applicationID string,
	) (ApplicationMailInfo, error)

	// Used by hermione - Candidacies related methods
	AddEmployerCandidacyComment(
		context.Context,
		employer.AddEmployerCandidacyCommentRequest,
	) (uuid.UUID, error)
	AddHubCandidacyComment(
		context.Context,
		hub.AddHubCandidacyCommentRequest,
	) (uuid.UUID, error)
	GetEmployerCandidacyComments(
		context.Context,
		common.GetCandidacyCommentsRequest,
	) ([]common.CandidacyComment, error)
	GetHubCandidacyComments(
		context.Context,
		common.GetCandidacyCommentsRequest,
	) ([]common.CandidacyComment, error)
	FilterEmployerCandidacyInfos(
		context.Context,
		employer.FilterCandidacyInfosRequest,
	) ([]employer.Candidacy, error)
	GetEmployerCandidacyInfo(
		context.Context,
		common.GetCandidacyInfoRequest,
	) (employer.Candidacy, error)
	GetHubCandidacyInfo(
		context.Context,
		common.GetCandidacyInfoRequest,
	) (hub.MyCandidacy, error)
	AddInterview(context.Context, AddInterviewRequest) error
	AddInterviewer(context.Context, AddInterviewerRequest) error
	RemoveInterviewer(context.Context, RemoveInterviewerRequest) error
	GetWatchersInfoByInterviewID(
		ctx context.Context,
		interviewID string,
	) (WatchersInfo, error)
	EmployerRSVPInterview(context.Context, common.RSVPInterviewRequest) error
	GetEmployerInterviewsByOpening(
		context.Context,
		employer.GetEmployerInterviewsByOpeningRequest,
	) ([]employer.EmployerInterview, error)
	GetEmployerInterviewsByCandidacy(
		context.Context,
		employer.GetEmployerInterviewsByCandidacyRequest,
	) ([]employer.EmployerInterview, error)
	GetInterview(context.Context, string) (employer.EmployerInterview, error)
	GetAssessment(
		context.Context,
		employer.GetAssessmentRequest,
	) (employer.Assessment, error)
	PutAssessment(
		context.Context,
		employer.PutAssessmentRequest,
	) error
	OfferToCandidate(context.Context, OfferToCandidateReq) error

	// Used by hermione - for Hub users
	AuthHubUser(c context.Context, token string) (HubUserTO, error)
	ChangeHubUserPassword(context.Context, uuid.UUID, string) error
	CreateApplication(context.Context, ApplyOpeningReq) error
	MyApplications(
		context.Context,
		hub.MyApplicationsRequest,
	) ([]hub.HubApplication, error)
	CreateHubUserToken(context.Context, HubTokenReq) error
	GetHubUserByTFACreds(context.Context, string, string) (HubUserTO, error)
	FindHubOpenings(
		context.Context,
		*hub.FindHubOpeningsRequest,
	) ([]hub.HubOpening, error)
	GetHubOpeningDetails(
		context.Context,
		hub.GetHubOpeningDetailsRequest,
	) (hub.HubOpeningDetails, error)
	GetMyCandidacies(
		context.Context,
		hub.MyCandidaciesRequest,
	) ([]hub.MyCandidacy, error)
	GetHubInterviewsByCandidacy(
		context.Context,
		hub.GetHubInterviewsByCandidacyRequest,
	) ([]hub.HubInterview, error)
	InitHubUserPasswordReset(context.Context, HubUserInitPasswordReset) error
	Logout(c context.Context, token string) error
	ResetHubUserPassword(context.Context, HubUserPasswordReset) error
	HubRSVPInterview(context.Context, hub.HubRSVPInterviewRequest) error
	GetCandidateInfo(context.Context, string) (CandidateInfo, error)

	// Opening tags
	FilterVTags(
		context.Context,
		common.FilterVTagsRequest,
	) ([]common.VTag, error)

	// Employers
	FilterEmployers(
		context.Context,
		hub.FilterEmployersRequest,
	) ([]hub.HubEmployer, error)

	// WorkHistory
	AddWorkHistory(
		context.Context,
		hub.AddWorkHistoryRequest,
	) (string, error)
	DeleteWorkHistory(context.Context, hub.DeleteWorkHistoryRequest) error
	ListWorkHistory(
		context.Context,
		hub.ListWorkHistoryRequest,
	) ([]hub.WorkHistory, error)
	UpdateWorkHistory(context.Context, hub.UpdateWorkHistoryRequest) error

	// Education
	AddEducation(context.Context, hub.AddEducationRequest) (string, error)
	DeleteEducation(context.Context, hub.DeleteEducationRequest) error
	ListEducation(
		context.Context,
		hub.ListEducationRequest,
	) ([]common.Education, error)
	ListHubUserEducation(
		context.Context,
		employer.ListHubUserEducationRequest,
	) ([]common.Education, error)
	FilterInstitutes(
		context.Context,
		hub.FilterInstitutesRequest,
	) ([]common.Institute, error)

	// Achievements
	AddAchievement(context.Context, hub.AddAchievementRequest) (string, error)
	ListAchievements(
		context.Context,
		hub.ListAchievementsRequest,
	) ([]common.Achievement, error)
	ListHubUserAchievements(
		context.Context,
		employer.ListHubUserAchievementsRequest,
	) ([]common.Achievement, error)
	DeleteAchievement(context.Context, hub.DeleteAchievementRequest) error

	// Used by hermione - Profile page related methods
	AddOfficialEmail(AddOfficialEmailReq) error
	GetMyOfficialEmails(context.Context) ([]hub.OfficialEmail, error)
	GetOfficialEmail(ctx context.Context, email string) (*OfficialEmail, error)
	UpdateOfficialEmailVerificationCode(
		ctx context.Context,
		req UpdateOfficialEmailVerificationCodeReq,
	) error
	VerifyOfficialEmail(ctx context.Context, email string, code string) error
	DeleteOfficialEmail(ctx context.Context, email string) error
	GetBio(ctx context.Context, handle string) (hub.Bio, error)
	GetEmployerViewBio(
		ctx context.Context,
		handle string,
	) (employer.EmployerViewBio, error)
	UpdateBio(ctx context.Context, bio hub.UpdateBioRequest) error
	UpdateProfilePictureWithCleanup(
		ctx context.Context,
		userID uuid.UUID,
		newPicturePath string,
	) error

	// Used by hermione - Colleagues related methods
	ConnectColleague(ctx context.Context, handle string) error
	GetMyColleagueApprovals(
		ctx context.Context,
		req hub.MyColleagueApprovalsRequest,
	) (hub.MyColleagueApprovals, error)
	GetMyColleagueSeeks(
		ctx context.Context,
		req hub.MyColleagueSeeksRequest,
	) (hub.MyColleagueSeeks, error)
	ApproveColleague(ctx context.Context, handle string) error
	RejectColleague(ctx context.Context, handle string) error
	UnlinkColleague(ctx context.Context, handle string) error
	FilterColleagues(
		ctx context.Context,
		req hub.FilterColleaguesRequest,
	) ([]hub.HubUserShort, error)
	GetMyEndorsementApprovals(
		ctx context.Context,
		req hub.MyEndorseApprovalsRequest,
	) (hub.MyEndorseApprovalsResponse, error)
	EndorseApplication(
		ctx context.Context,
		endorseReq hub.EndorseApplicationRequest,
	) error
	RejectEndorsement(
		ctx context.Context,
		rejectReq hub.RejectEndorsementRequest,
	) error

	// Used by hermione - HubUsers related methods
	GetHubUsersByHandles(
		ctx context.Context,
		handles []common.Handle,
	) ([]HubUserContact, error)
	InviteHubUser(ctx context.Context, inviteHubUserReq InviteHubUserReq) error
	OnboardHubUser(
		context.Context,
		OnboardHubUserReq,
	) (generatedHandle string, err error)
	CheckHandleAvailability(
		ctx context.Context,
		handle common.Handle,
	) (hub.CheckHandleAvailabilityResponse, error)

	// Used by granger
	PruneOfficialEmailCodes(ctx context.Context) error
	GetStaleFiles(ctx context.Context, limit int) ([]StaleFile, error)
	MarkFileCleaned(
		ctx context.Context,
		fileID uuid.UUID,
		cleanedAt time.Time,
	) error

	GetUnscoredApplication(
		ctx context.Context,
		limit int,
	) (*UnscoredApplicationBatch, error)
	SaveApplicationScores(ctx context.Context, scores []ApplicationScore) error

	// Employer settings
	ChangeCoolOffPeriod(ctx context.Context, coolOffPeriod int32) error
	GetCoolOffPeriod(ctx context.Context) (int32, error)

	// Used by hermione - Posts related methods
	AddPost(req AddPostRequest) error
	GetUserPosts(
		ctx context.Context,
		req hub.GetUserPostsRequest,
	) (hub.GetUserPostsResponse, error)
}
