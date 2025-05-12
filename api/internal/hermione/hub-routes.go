package hermione

import (
	"net/http"

	ach "github.com/vetchium/vetchium/api/internal/hermione/achievements"
	app "github.com/vetchium/vetchium/api/internal/hermione/applications"
	ca "github.com/vetchium/vetchium/api/internal/hermione/candidacy"
	co "github.com/vetchium/vetchium/api/internal/hermione/colleagues"
	ed "github.com/vetchium/vetchium/api/internal/hermione/education"
	ha "github.com/vetchium/vetchium/api/internal/hermione/hubauth"
	he "github.com/vetchium/vetchium/api/internal/hermione/hubemp"
	ho "github.com/vetchium/vetchium/api/internal/hermione/hubopenings"
	hu "github.com/vetchium/vetchium/api/internal/hermione/hubusers"
	in "github.com/vetchium/vetchium/api/internal/hermione/interview"
	po "github.com/vetchium/vetchium/api/internal/hermione/posts"
	pp "github.com/vetchium/vetchium/api/internal/hermione/profilepage"
	wh "github.com/vetchium/vetchium/api/internal/hermione/workhistory"
	"github.com/vetchium/vetchium/typespec/hub"
)

func RegisterHubRoutes(h *Hermione) {
	// Unprotected routes
	http.HandleFunc("/hub/login", ha.Login(h))
	http.HandleFunc("/hub/tfa", ha.HubTFA(h))
	http.HandleFunc("/hub/logout", ha.Logout(h))
	http.HandleFunc("/hub/forgot-password", ha.ForgotPassword(h))
	http.HandleFunc("/hub/reset-password", ha.ResetPassword(h))
	http.HandleFunc("/hub/onboard-user", hu.OnboardHubUser(h))
	http.HandleFunc("/hub/signup", hu.SignupHubUser(h))

	h.mw.Guard(
		"/hub/change-email-address",
		hu.ChangeEmailAddress(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/get-my-handle",
		ha.GetMyHandle(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/change-password",
		ha.ChangePassword(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/invite-hub-user",
		hu.InviteHubUser(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/check-handle-availability",
		hu.CheckHandleAvailability(h),
		[]hub.HubUserTier{hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/set-handle",
		hu.SetHandle(h),
		[]hub.HubUserTier{hub.PaidHubUserTier},
	)

	// Official Email related endpoints
	h.mw.Guard(
		"/hub/add-official-email",
		pp.AddOfficialEmail(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/verify-official-email",
		pp.VerifyOfficialEmail(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/trigger-verification",
		pp.TriggerVerification(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/delete-official-email",
		pp.DeleteOfficialEmail(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-official-emails",
		pp.MyOfficialEmails(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	// ProfilePage related endpoints
	h.mw.Guard(
		"/hub/get-bio",
		pp.GetBio(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/update-bio",
		pp.UpdateBio(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/upload-profile-picture",
		pp.UploadProfilePicture(h),
		[]hub.HubUserTier{hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/profile-picture/",
		pp.GetProfilePicture(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/remove-profile-picture",
		pp.RemoveProfilePicture(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-tier",
		pp.MyTier(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/add-education",
		ed.AddEducation(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/filter-institutes",
		ed.FilterInstitutes(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/delete-education",
		ed.DeleteEducation(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/list-education",
		ed.ListEducation(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/add-achievement",
		ach.AddAchievement(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/list-achievements",
		ach.ListAchievements(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/delete-achievement",
		ach.DeleteAchievement(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/find-openings",
		ho.FindHubOpenings(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/filter-vtags",
		he.FilterVTags(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-opening-details",
		ho.GetOpeningDetails(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/apply-for-opening",
		ho.ApplyForOpening(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-applications",
		app.MyApplications(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/withdraw-application",
		app.WithdrawApplication(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/add-candidacy-comment",
		ca.HubAddComment(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-candidacy-comments",
		ca.HubGetComments(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-my-candidacies",
		ca.MyCandidacies(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-candidacy-info",
		ca.GetHubCandidacyInfo(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-interviews-by-candidacy",
		in.GetHubInterviewsByCandidacy(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/rsvp-interview",
		in.HubRSVPInterview(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/filter-employers",
		he.FilterEmployers(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	// WorkHistory related endpoints
	h.mw.Guard(
		"/hub/add-work-history",
		wh.AddWorkHistory(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/delete-work-history",
		wh.DeleteWorkHistory(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/list-work-history",
		wh.ListWorkHistory(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/update-work-history",
		wh.UpdateWorkHistory(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	// Colleague related endpoints
	h.mw.Guard(
		"/hub/connect-colleague",
		co.ConnectColleague(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/unlink-colleague",
		co.UnlinkColleague(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-colleague-approvals",
		co.MyColleagueApprovals(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-endorse-approvals",
		co.MyEndorseApprovals(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/my-colleague-seeks",
		co.MyColleagueSeeks(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/endorse-application",
		co.EndorseApplication(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/reject-endorsement",
		co.RejectEndorsement(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/approve-colleague",
		co.ApproveColleague(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/reject-colleague",
		co.RejectColleague(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/filter-colleagues",
		co.FilterColleagues(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)

	h.mw.Guard(
		"/hub/add-post",
		po.AddPost(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-user-posts",
		po.GetUserPosts(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/follow-user",
		po.FollowUser(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/unfollow-user",
		po.UnfollowUser(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-follow-status",
		po.GetFollowStatus(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-my-home-timeline",
		po.GetMyHomeTimeline(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/get-post-details",
		po.GetPostDetails(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/upvote-user-post",
		po.UpvoteUserPost(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/downvote-user-post",
		po.DownvoteUserPost(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
	h.mw.Guard(
		"/hub/unvote-user-post",
		po.UnvoteUserPost(h),
		[]hub.HubUserTier{hub.FreeHubUserTier, hub.PaidHubUserTier},
	)
}
