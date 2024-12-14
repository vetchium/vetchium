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

func AddInterviewers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddInterviewers")
		var addInterviewersReq employer.AddInterviewersRequest
		err := json.NewDecoder(r.Body).Decode(&addInterviewersReq)
		if err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addInterviewersReq) {
			h.Dbg("validation failed", "addInterviewersReq", addInterviewersReq)
			return
		}
		h.Dbg("validated", "addInterviewersReq", addInterviewersReq)

		notificationEmail, err := h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.AddInterviewers,
				Args: map[string]string{
					"InterviewURL": "TODO",
				},
				EmailFrom: vetchi.EmailFrom,
				EmailTo:   addInterviewersReq.OrgUserEmails,

				// TODO: This should be dynamic and come from hedwig
				Subject: "Added as an Interviewer",
			})
		if err != nil {
			h.Dbg("generating email failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().AddInterviewers(r.Context(), db.AddInterviewersRequest{
			InterviewID: addInterviewersReq.InterviewID,
			Email:       notificationEmail,
		})
		if err != nil {
			h.Dbg("adding interviewers failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
