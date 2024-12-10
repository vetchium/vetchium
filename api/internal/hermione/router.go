package hermione

import (
	"fmt"
	"net/http"

	app "github.com/psankar/vetchi/api/internal/hermione/applications"
	"github.com/psankar/vetchi/api/internal/hermione/costcenter"
	ea "github.com/psankar/vetchi/api/internal/hermione/employerauth"
	ha "github.com/psankar/vetchi/api/internal/hermione/hubauth"
	ho "github.com/psankar/vetchi/api/internal/hermione/hubopenings"
	"github.com/psankar/vetchi/api/internal/hermione/locations"
	"github.com/psankar/vetchi/api/internal/hermione/openings"
	"github.com/psankar/vetchi/api/internal/hermione/orgusers"
	"github.com/psankar/vetchi/typespec/common"
)

func (h *Hermione) Run() error {
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
	http.HandleFunc("/employer/signup-org-user", orgusers.SignupOrgUser(h))

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

	wrap := func(fn http.Handler) http.Handler {
		return h.mw.HubWrap(fn)
	}
	// Hub related endpoints
	http.HandleFunc("/hub/login", ha.Login(h))
	http.HandleFunc("/hub/tfa", ha.HubTFA(h))
	http.Handle("/hub/get-my-handle", wrap(ha.GetMyHandle(h)))
	http.HandleFunc("/hub/logout", ha.Logout(h))

	http.HandleFunc("/hub/forgot-password", ha.ForgotPassword(h))
	http.HandleFunc("/hub/reset-password", ha.ResetPassword(h))
	http.Handle("/hub/change-password", wrap(ha.ChangePassword(h)))

	http.Handle("/hub/find-openings", wrap(ho.FindHubOpenings(h)))
	http.Handle("/hub/apply-for-opening", wrap(ho.ApplyForOpening(h)))
	http.Handle("/hub/my-applications", wrap(app.MyApplications(h)))
	http.Handle("/hub/withdraw-application", wrap(app.WithdrawApplication(h)))

	port := fmt.Sprintf(":%d", h.Config().Port)
	return http.ListenAndServe(port, nil)

	/*

		TODO:

		   - /employer/add-candidacy-comment
		   - /employer/get-candidacies-info
		   - /hub/add-candidacy-comment
		   - /hub/get-my-candidacies

		   - /employer/add-interview
		   - /employer/get-interviews

	*/

}
