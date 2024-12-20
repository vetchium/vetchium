package interview

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/employer"
)

func AddInterviewer(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddInterviewers")
		var addInterviewerReq employer.AddInterviewerRequest
		err := json.NewDecoder(r.Body).Decode(&addInterviewerReq)
		if err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addInterviewerReq) {
			h.Dbg("validation failed", "addInterviewerReq", addInterviewerReq)
			return
		}
		h.Dbg("validated", "addInterviewerReq", addInterviewerReq)

		ctx := r.Context()

		orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		interviewer, err := h.DB().
			GetOrgUserByEmail(ctx, addInterviewerReq.OrgUserEmail)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get org user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if interviewer.OrgUserState != employer.ActiveOrgUserState &&
			interviewer.OrgUserState != employer.AddedOrgUserState &&
			interviewer.OrgUserState != employer.ReplicatedOrgUserState {
			h.Err("interviewer is not active", "interviewer", interviewer)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		watchersInfo, err := h.DB().
			GetWatchersInfoByInterviewID(ctx, addInterviewerReq.InterviewID)
		if err != nil {
			h.Err("failed to get watchers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		watcherEmailMap := make(map[string]struct{})
		watcherEmailMap[addInterviewerReq.OrgUserEmail] = struct{}{}
		watcherEmailMap[watchersInfo.HiringManager.Email] = struct{}{}
		watcherEmailMap[watchersInfo.Recruiter.Email] = struct{}{}

		// Should we notify the Applicant also ? Perhaps not, let us not
		// distract them with too many emails, when they should be preparing
		// for the interview
		for _, watcher := range watchersInfo.Watchers {
			watcherEmailMap[watcher.Email] = struct{}{}
		}

		watcherEmailRecipients := make([]string, 0, len(watcherEmailMap))
		for email := range watcherEmailMap {
			watcherEmailRecipients = append(watcherEmailRecipients, email)
		}

		watcherNotification, err := h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.NotifyNewInterviewer,
				Args: map[string]string{
					"InterviewerName": interviewer.Name,
					"InterviewURL":    "TODO",
				},
				EmailFrom: vetchi.EmailFrom,
				EmailTo:   watcherEmailRecipients,

				// TODO: This should be dynamic and come from hedwig
				Subject: "New Interviewer Added",
			})
		if err != nil {
			h.Dbg("generating email failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		interviewerNotification, err := h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.NotifyNewInterviewer,
				Args: map[string]string{
					"InterviewURL": "TODO",
				},
				EmailFrom: vetchi.EmailFrom,
				EmailTo:   []string{addInterviewerReq.OrgUserEmail},

				// TODO: This should be dynamic and come from hedwig
				Subject: "Added as an Interviewer",
			})
		if err != nil {
			h.Dbg("generating email failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		candidacyComment := fmt.Sprintf(
			"%s added %s as an interviewer for %s TODO: i18n",
			orgUser.Name,
			interviewer.Name,
			addInterviewerReq.InterviewID,
		)

		err = h.DB().AddInterviewer(ctx, db.AddInterviewerRequest{
			InterviewID:                  addInterviewerReq.InterviewID,
			InterviewerEmailAddr:         addInterviewerReq.OrgUserEmail,
			InterviewerNotificationEmail: interviewerNotification,
			WatcherNotificationEmail:     watcherNotification,
			CandidacyComment:             candidacyComment,
		})
		if err != nil {
			// TODO: The error codes needs to be standardized across *.tsp files
			// and all the endpoints
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			if errors.Is(err, db.ErrInterviewerNotActive) ||
				errors.Is(err, db.ErrInvalidInterviewState) {
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			if errors.Is(err, db.ErrNoInterview) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("adding interviewers failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
