package interview

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

func AddInterview(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddInterview")
		var addInterviewReq employer.AddInterviewRequest
		if err := json.NewDecoder(r.Body).Decode(&addInterviewReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addInterviewReq) {
			h.Dbg("validation failed", "addInterviewReq", addInterviewReq)
			return
		}
		h.Dbg("validated", "addInterviewReq", addInterviewReq)
		h.Dbg("timestamp debug",
			"start_time_utc", addInterviewReq.StartTime.UTC(),
			"end_time_utc", addInterviewReq.EndTime.UTC(),
			"start_time_local", addInterviewReq.StartTime.Local(),
			"end_time_local", addInterviewReq.EndTime.Local(),
		)

		if len(addInterviewReq.InterviewerEmails) > 5 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(common.ValidationErrors{
				Errors: []string{"interviewer_emails"},
			})
			return
		}

		interviewID := util.RandomUniqueID(vetchi.InterviewIDLenBytes)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Prepare email notifications
		var interviewerNotification, watcherNotification, applicantNotification db.Email
		var err error

		// Prepare applicant notification
		emailArgs := map[string]string{
			"InterviewURL":    "TODO",
			"InterviewType":   string(addInterviewReq.InterviewType),
			"StartTime":       addInterviewReq.StartTime.String(),
			"EndTime":         addInterviewReq.EndTime.String(),
			"Description":     addInterviewReq.Description,
			"InterviewerName": orgUser.Name,
		}

		if len(addInterviewReq.InterviewerEmails) > 0 {
			// We'll get the actual names from the database during the transaction
			emailArgs["Interviewers"] = "Your interviewers will be: " +
				"[Interviewer names will be added during the transaction]"
		}

		applicantNotification, err = h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.NotifyApplicantInterview,
				Args:         emailArgs,
				EmailFrom:    vetchi.EmailFrom,
				Subject:      "Interview Scheduled",
			})
		if err != nil {
			h.Dbg("generating applicant email failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if len(addInterviewReq.InterviewerEmails) > 0 {
			interviewerNotification, err = h.Hedwig().
				GenerateEmail(hedwig.GenerateEmailReq{
					TemplateName: hedwig.NotifyNewInterviewer,
					Args: map[string]string{
						"InterviewURL": "TODO",
					},
					EmailFrom: vetchi.EmailFrom,
					EmailTo:   addInterviewReq.InterviewerEmails,
					Subject:   "Added as an Interviewer",
				})
			if err != nil {
				h.Dbg("generating interviewer email failed", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			watcherNotification, err = h.Hedwig().
				GenerateEmail(hedwig.GenerateEmailReq{
					TemplateName: hedwig.NotifyWatchersNewInterviewer,
					Args: map[string]string{
						"InterviewURL": "TODO",
					},
					EmailFrom: vetchi.EmailFrom,
					EmailTo: []string{
						orgUser.Email,
					}, // TODO: Add other watchers
					Subject: "New Interviewer Added",
				})
			if err != nil {
				h.Dbg("generating watcher email failed", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		candidacyComment := fmt.Sprintf(
			"%s scheduled a new interview of type %s TODO: i18n",
			orgUser.Name,
			addInterviewReq.InterviewType,
		)

		err = h.DB().AddInterview(r.Context(), db.AddInterviewRequest{
			AddInterviewRequest:          addInterviewReq,
			InterviewID:                  interviewID,
			InterviewerNotificationEmail: interviewerNotification,
			WatcherNotificationEmail:     watcherNotification,
			ApplicantNotificationEmail:   applicantNotification,
			CandidacyComment:             candidacyComment,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoCandidacy) {
				h.Dbg("no candidacy found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrInvalidCandidacyState) {
				h.Dbg("candidacy not in valid state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to add interview", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added interview", "interviewID", interviewID)
		err = json.NewEncoder(w).Encode(employer.AddInterviewResponse{
			InterviewID: interviewID,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
