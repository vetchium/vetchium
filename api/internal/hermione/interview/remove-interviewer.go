package interview

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/employer"
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

		ctx := r.Context()

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

		orgUser, ok := ctx.Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		removedInterviewer, err := h.DB().
			GetOrgUserByEmail(ctx, removeInterviewerReq.OrgUserEmail)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				h.Err("failed to get removed interviewer name", "err", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Err("failed to get removed interviewer name", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// No need to check if the removed interviewer is active

		candidacyComment := fmt.Sprintf(
			"%s removed %s as an interviewer",
			orgUser.Name,
			removedInterviewer.Name,
		)

		err = h.DB().RemoveInterviewer(ctx, db.RemoveInterviewerRequest{
			InterviewID:                         removeInterviewerReq.InterviewID,
			RemovedInterviewerEmailAddr:         removeInterviewerReq.OrgUserEmail,
			RemovedInterviewerEmailNotification: removedInterviewerNotify,
			CandidacyComment:                    candidacyComment,
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
