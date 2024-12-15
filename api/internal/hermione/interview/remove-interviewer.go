package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/employer"
)

func RemoveInterviewer(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveInterviewer")
		var removeInterviewerReq employer.RemoveInterviewerRequest
		err := json.NewDecoder(r.Body).Decode(&removeInterviewerReq)
		if err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &removeInterviewerReq) {
			h.Dbg("Validation failed", "req", removeInterviewerReq)
			return
		}
		h.Dbg("validated", "removeInterviewerReq", removeInterviewerReq)

		removedInterviewerNotify, err := h.Hedwig().
			GenerateEmail(hedwig.GenerateEmailReq{
				TemplateName: hedwig.RemovedInterviewerNotify,
				Args: map[string]string{
					"InterviewURL": "TODO",
				},
				EmailFrom: vetchi.EmailFrom,
				EmailTo:   []string{removeInterviewerReq.OrgUserEmail},

				// TODO: This should be dynamic and come from hedwig
				Subject: "Your participation as an interviewer is removed",
			})
		if err != nil {
			h.Dbg("Error generating email", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().RemoveInterviewer(r.Context(), db.RemoveInterviewerRequest{
			InterviewID:                         removeInterviewerReq.InterviewID,
			RemovedInterviewerEmailAddr:         removeInterviewerReq.OrgUserEmail,
			RemovedInterviewerEmailNotification: removedInterviewerNotify,
			CandidacyComment:                    "TODO: i18n TODO: orgUserName1 removed orgUserName2 as an interviewer",
		})
		if err != nil {
			if errors.Is(err, db.ErrInvalidInterviewState) {
				h.Dbg("Invalid interview state", "err", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("Error removing interviewer", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}
