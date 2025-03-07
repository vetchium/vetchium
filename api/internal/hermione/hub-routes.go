package hermione

import (
	"net/http"

	app "github.com/psankar/vetchi/api/internal/hermione/applications"
	ca "github.com/psankar/vetchi/api/internal/hermione/candidacy"
	co "github.com/psankar/vetchi/api/internal/hermione/colleagues"
	ha "github.com/psankar/vetchi/api/internal/hermione/hubauth"
	he "github.com/psankar/vetchi/api/internal/hermione/hubemp"
	ho "github.com/psankar/vetchi/api/internal/hermione/hubopenings"
	hu "github.com/psankar/vetchi/api/internal/hermione/hubusers"
	in "github.com/psankar/vetchi/api/internal/hermione/interview"
	pp "github.com/psankar/vetchi/api/internal/hermione/profilepage"
	wh "github.com/psankar/vetchi/api/internal/hermione/workhistory"
)

func RegisterHubRoutes(h *Hermione) {
	wrap := func(fn http.Handler) http.Handler {
		return h.mw.HubWrap(fn)
	}

	http.HandleFunc("/hub/login", ha.Login(h))
	http.HandleFunc("/hub/tfa", ha.HubTFA(h))
	http.Handle("/hub/get-my-handle", wrap(ha.GetMyHandle(h)))
	http.HandleFunc("/hub/logout", ha.Logout(h))

	http.HandleFunc("/hub/forgot-password", ha.ForgotPassword(h))
	http.HandleFunc("/hub/reset-password", ha.ResetPassword(h))
	http.Handle("/hub/change-password", wrap(ha.ChangePassword(h)))

	http.Handle("/hub/invite-hub-user", wrap(hu.InviteHubUser(h)))
	http.Handle("/hub/onboard-user", hu.OnboardHubUser(h))
	http.Handle(
		"/hub/check-handle-availability",
		wrap(hu.CheckHandleAvailability(h)),
	)
	http.Handle("/hub/set-handle", wrap(hu.SetHandle(h)))

	// Official Email related endpoints
	http.Handle("/hub/add-official-email", wrap(pp.AddOfficialEmail(h)))
	http.Handle("/hub/verify-official-email", wrap(pp.VerifyOfficialEmail(h)))
	http.Handle("/hub/trigger-verification", wrap(pp.TriggerVerification(h)))
	http.Handle("/hub/delete-official-email", wrap(pp.DeleteOfficialEmail(h)))
	http.Handle("/hub/my-official-emails", wrap(pp.MyOfficialEmails(h)))

	// ProfilePage related endpoints
	http.Handle("/hub/get-bio", wrap(pp.GetBio(h)))
	http.Handle("/hub/update-bio", wrap(pp.UpdateBio(h)))
	http.Handle("/hub/upload-profile-picture", wrap(pp.UploadProfilePicture(h)))
	http.Handle("/hub/profile-picture/", wrap(pp.GetProfilePicture(h)))
	http.Handle("/hub/remove-profile-picture", wrap(pp.RemoveProfilePicture(h)))

	http.Handle("/hub/find-openings", wrap(ho.FindHubOpenings(h)))
	http.Handle("/hub/filter-opening-tags", wrap(he.FilterOpeningTags(h)))
	http.Handle("/hub/get-opening-details", wrap(ho.GetOpeningDetails(h)))
	http.Handle("/hub/apply-for-opening", wrap(ho.ApplyForOpening(h)))
	http.Handle("/hub/my-applications", wrap(app.MyApplications(h)))
	http.Handle("/hub/withdraw-application", wrap(app.WithdrawApplication(h)))
	http.Handle("/hub/add-candidacy-comment", wrap(ca.HubAddComment(h)))
	http.Handle("/hub/get-candidacy-comments", wrap(ca.HubGetComments(h)))
	http.Handle("/hub/get-my-candidacies", wrap(ca.MyCandidacies(h)))
	http.Handle("/hub/get-candidacy-info", wrap(ca.GetHubCandidacyInfo(h)))
	http.Handle(
		"/hub/get-interviews-by-candidacy",
		wrap(in.GetHubInterviewsByCandidacy(h)),
	)
	http.Handle("/hub/rsvp-interview", wrap(in.HubRSVPInterview(h)))
	http.Handle("/hub/filter-employers", wrap(he.FilterEmployers(h)))

	// WorkHistory related endpoints
	http.Handle("/hub/add-work-history", wrap(wh.AddWorkHistory(h)))
	http.Handle("/hub/delete-work-history", wrap(wh.DeleteWorkHistory(h)))
	http.Handle("/hub/list-work-history", wrap(wh.ListWorkHistory(h)))
	http.Handle("/hub/update-work-history", wrap(wh.UpdateWorkHistory(h)))

	// Colleague related endpoints
	http.Handle("/hub/connect-colleague", wrap(co.ConnectColleague(h)))
	http.Handle("/hub/unlink-colleague", wrap(co.UnlinkColleague(h)))
	http.Handle("/hub/my-colleague-approvals", wrap(co.MyColleagueApprovals(h)))
	http.Handle("/hub/my-endorse-approvals", wrap(co.MyEndorseApprovals(h)))
	http.Handle("/hub/my-colleague-seeks", wrap(co.MyColleagueSeeks(h)))
	http.Handle("/hub/endorse-application", wrap(co.EndorseApplication(h)))
	http.Handle("/hub/reject-endorsement", wrap(co.RejectEndorsement(h)))
	http.Handle("/hub/approve-colleague", wrap(co.ApproveColleague(h)))
	http.Handle("/hub/reject-colleague", wrap(co.RejectColleague(h)))
	http.Handle("/hub/filter-colleagues", wrap(co.FilterColleagues(h)))
}
