package hermione

import (
	"net/http"

	"github.com/psankar/vetchi/api/internal/hermione/costcenter"
	ea "github.com/psankar/vetchi/api/internal/hermione/employerauth"
	"github.com/psankar/vetchi/api/internal/hermione/locations"
	"github.com/psankar/vetchi/api/internal/hermione/openings"
	"github.com/psankar/vetchi/api/internal/hermione/orgusers"
	"github.com/psankar/vetchi/api/pkg/vetchi"
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
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-centers",
		costcenter.GetCostCenters(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.CostCentersCRUD,
			vetchi.CostCentersViewer,
		},
	)
	h.mw.Protect(
		"/employer/defunct-cost-center",
		costcenter.DefunctCostCenter(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/rename-cost-center",
		costcenter.RenameCostCenter(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/update-cost-center",
		costcenter.UpdateCostCenter(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersCRUD},
	)
	h.mw.Protect(
		"/employer/get-cost-center",
		costcenter.GetCostCenter(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.CostCentersViewer},
	)

	// Location related endpoints
	h.mw.Protect(
		"/employer/add-location",
		locations.AddLocation(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/defunct-location",
		locations.DefunctLocation(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/get-locations",
		locations.GetLocations(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/get-location",
		locations.GetLocation(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.LocationsCRUD,
			vetchi.LocationsViewer,
		},
	)
	h.mw.Protect(
		"/employer/rename-location",
		locations.RenameLocation(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)
	h.mw.Protect(
		"/employer/update-location",
		locations.UpdateLocation(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.LocationsCRUD},
	)

	// OrgUser related endpoints
	h.mw.Protect(
		"/employer/add-org-user",
		orgusers.AddOrgUser(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/update-org-user",
		orgusers.UpdateOrgUser(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/disable-org-user",
		orgusers.DisableOrgUser(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/enable-org-user",
		orgusers.EnableOrgUser(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OrgUsersCRUD},
	)
	h.mw.Protect(
		"/employer/filter-org-users",
		orgusers.FilterOrgUsers(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.OrgUsersCRUD,
			vetchi.OrgUsersViewer,
		},
	)
	http.HandleFunc("/employer/signup-org-user", orgusers.SignupOrgUser(h))

	// Opening related endpoints
	h.mw.Protect(
		"/employer/create-opening",
		openings.CreateOpening(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening",
		openings.GetOpening(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.OpeningsCRUD,
			vetchi.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/filter-openings",
		openings.FilterOpenings(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.OpeningsCRUD,
			vetchi.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/update-opening",
		openings.UpdateOpening(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/get-opening-watchers",
		openings.GetOpeningWatchers(h),
		[]vetchi.OrgUserRole{
			vetchi.Admin,
			vetchi.OpeningsCRUD,
			vetchi.OpeningsViewer,
		},
	)
	h.mw.Protect(
		"/employer/add-opening-watchers",
		openings.AddOpeningWatchers(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/remove-opening-watcher",
		openings.RemoveOpeningWatcher(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OpeningsCRUD},
	)
	h.mw.Protect(
		"/employer/change-opening-state",
		openings.ChangeOpeningState(h),
		[]vetchi.OrgUserRole{vetchi.Admin, vetchi.OpeningsCRUD},
	)

	return http.ListenAndServe(h.port, nil)
}
