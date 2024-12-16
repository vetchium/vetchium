package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
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

		notificationEmail, err := h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.AddInterviewers,
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

		err = h.DB().AddInterviewer(r.Context(), db.AddInterviewerRequest{
			InterviewID:                  addInterviewerReq.InterviewID,
			InterviewerEmailAddr:         addInterviewerReq.OrgUserEmail,
			InterviewerNotificationEmail: notificationEmail,
		})
		if err != nil {
			h.Dbg("adding interviewers failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
