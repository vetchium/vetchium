package hermione

import (
	"net/http"

	"github.com/vetchium/vetchium/api/internal/hermione/achievements"
	app "github.com/vetchium/vetchium/api/internal/hermione/applications"
	"github.com/vetchium/vetchium/api/internal/hermione/candidacy"
	"github.com/vetchium/vetchium/api/internal/hermione/costcenter"
	"github.com/vetchium/vetchium/api/internal/hermione/education"
	ea "github.com/vetchium/vetchium/api/internal/hermione/employerauth"
	"github.com/vetchium/vetchium/api/internal/hermione/employersettings"
	he "github.com/vetchium/vetchium/api/internal/hermione/hubemp"
	"github.com/vetchium/vetchium/api/internal/hermione/interview"
	"github.com/vetchium/vetchium/api/internal/hermione/locations"
	"github.com/vetchium/vetchium/api/internal/hermione/openings"
	"github.com/vetchium/vetchium/api/internal/hermione/orgusers"
	pp "github.com/vetchium/vetchium/api/internal/hermione/profilepage"
	"github.com/vetchium/vetchium/typespec/common"
)

func RegisterEmployerRoutes(h *Hermione) {
	// Authentication related endpoints
	http.HandleFunc("/employer/get-onboard-status", ea.GetOnboardStatus(h))
	http.HandleFunc("/employer/set-onboard-password", ea.SetOnboardPassword(h))
	http.HandleFunc("/employer/signin", ea.EmployerSignin(h))
	http.HandleFunc("/employer/tfa", ea.EmployerTFA(h))

	// CostCenter related endpoints
	h.mw.Protect(
		"/employer/add-cost-center",
		costcenter.AddCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-centers",
		costcenter.GetCostCenters(h),
		[]common.OrgUserRole{
			common.Admin,
			common.CostCentersCRUD,
			common.CostCentersViewer,
		},
	)
	h.mw.Protect(
		"/employer/defunct-cost-center",
		costcenter.DefunctCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/rename-cost-center",
		costcenter.RenameCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/update-cost-center",
		costcenter.UpdateCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-center",
		costcenter.GetCostCenter(h),
		[]common.OrgUserRole{common.Admin, common.CostCentersViewer},
	)

	// Location related endpoints
	h.mw.Protect(
		"/employer/add-location",
		locations.AddLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/defunct-location",
		locations.DefunctLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/get-locations",
		locations.GetLocations(h),
		[]common.OrgUserRole{
			common.Admin,
			common.LocationsCRUD,
			common.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/get-location",
		locations.GetLocation(h),
		[]common.OrgUserRole{
			common.Admin,
			common.LocationsCRUD,
			common.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/rename-location",
		locations.RenameLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/update-location",
		locations.UpdateLocation(h),
		[]common.OrgUserRole{common.Admin, common.LocationsCRUD},
	)

	// OrgUser related endpoints
	h.mw.Protect(
		"/employer/add-org-user",
		orgusers.AddOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/update-org-user",
		orgusers.UpdateOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/disable-org-user",
		orgusers.DisableOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/enable-org-user",
		orgusers.EnableOrgUser(h),
		[]common.OrgUserRole{common.Admin, common.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/filter-org-users",
		orgusers.FilterOrgUsers(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OrgUsersCRUD,
			common.OrgUsersViewer,
		},
	)
	http.HandleFunc("/employer/signup-orguser", orgusers.SignupOrgUser(h))

	// Opening related endpoints
	h.mw.Protect(
		"/employer/create-opening",
		openings.CreateOpening(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening",
		openings.GetOpening(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/filter-openings",
		openings.FilterOpenings(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/update-opening",
		openings.UpdateOpening(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening-watchers",
		openings.GetOpeningWatchers(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/add-opening-watchers",
		openings.AddOpeningWatchers(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/remove-opening-watcher",
		openings.RemoveOpeningWatcher(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/change-opening-state",
		openings.ChangeOpeningState(h),
		[]common.OrgUserRole{common.Admin, common.OpeningsCRUD},
	)

	// Opening tags related endpoints
	h.mw.Protect(
		"/employer/filter-opening-tags",
		he.FilterVTags(h),
		[]common.OrgUserRole{
			common.Admin,
			common.OpeningsCRUD,
			common.OpeningsViewer,
		},
	)

	// Application related endpoints
	h.mw.Protect(
		"/employer/get-applications",
		app.GetApplications(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)

	h.mw.Protect(
		"/employer/get-resume",
		app.GetResume(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)

	h.mw.Protect(
		"/employer/set-application-color-tag",
		app.SetApplicationColorTag(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/remove-application-color-tag",
		app.RemoveApplicationColorTag(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/shortlist-application",
		app.ShortlistApplication(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	h.mw.Protect(
		"/employer/reject-application",
		app.RejectApplication(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
		},
	)

	// Used by employer - Candidacies
	h.mw.Protect(
		"/employer/add-candidacy-comment",
		candidacy.EmployerAddComment(h),
		[]common.OrgUserRole{common.AnyOrgUser},
	)

	h.mw.Protect(
		"/employer/get-candidacy-comments",
		candidacy.EmployerGetComments(h),
		[]common.OrgUserRole{common.AnyOrgUser},
	)

	h.mw.Protect(
		"/employer/filter-candidacy-infos",
		candidacy.FilterCandidacyInfos(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.AnyOrgUser},
	)

	h.mw.Protect(
		"/employer/get-candidacy-info",
		candidacy.GetEmployerCandidacyInfo(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/offer-to-candidate",
		candidacy.OfferToCandidate(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	// Used by employer - Interviews
	h.mw.Protect(
		"/employer/add-interview",
		interview.AddInterview(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)
	h.mw.Protect(
		"/employer/add-interviewer",
		interview.AddInterviewer(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	h.mw.Protect(
		"/employer/remove-interviewer",
		interview.RemoveInterviewer(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	h.mw.Protect(
		"/employer/rsvp-interview",
		interview.EmployerRSVPInterview(h),
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/get-interviews-by-opening",
		interview.GetInterviewsByOpening(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/get-interviews-by-candidacy",
		interview.GetEmployerInterviewsByCandidacy(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/get-assessment",
		interview.EmployerGetAssessment(h),
		// TODO: It is unclear what roles should be required here
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/get-interview-details",
		interview.GetInterviewDetails(h),
		[]common.OrgUserRole{common.AnyOrgUser},
	)
	h.mw.Protect(
		"/employer/put-assessment",
		interview.EmployerPutAssessment(h),
		[]common.OrgUserRole{common.AnyOrgUser},
	)

	// Hub user profile related endpoints for employer
	h.mw.Protect(
		"/employer/get-hub-user-bio",
		pp.GetHubUserBio(h),
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/get-hub-user-profile-picture/",
		func(w http.ResponseWriter, r *http.Request) {
			http.StripPrefix(
				"/employer/get-hub-user-profile-picture/",
				pp.GetHubUserProfilePicture(h),
			).ServeHTTP(w, r)
		},
		[]common.OrgUserRole{
			common.Admin,
			common.ApplicationsCRUD,
			common.ApplicationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/list-hub-user-education",
		education.ListHubUserEducation(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	h.mw.Protect(
		"/employer/list-hub-user-achievements",
		achievements.ListHubUserAchievements(h),
		[]common.OrgUserRole{common.Admin, common.ApplicationsCRUD},
	)

	// settings related endpoints
	h.mw.Protect(
		"/employer/change-cool-off-period",
		employersettings.ChangeCoolOffPeriod(h),
		[]common.OrgUserRole{common.Admin},
	)

	h.mw.Protect(
		"/employer/get-cool-off-period",
		employersettings.GetCoolOffPeriod(h),
		[]common.OrgUserRole{common.Admin},
	)
}
