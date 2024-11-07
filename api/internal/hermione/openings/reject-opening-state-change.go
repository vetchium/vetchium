package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func RejectOpeningStateChange(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RejectOpeningStateChange")
		var rejectReq vetchi.RejectOpeningStateChangeRequest
		err := json.NewDecoder(r.Body).Decode(&rejectReq)
		if err != nil {
			h.Dbg("failed to decode reject state change request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &rejectReq) {
			h.Dbg("validation failed", "rejectReq", rejectReq)
			return
		}
		h.Dbg("validated", "rejectReq", rejectReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().
			RejectOpeningStateChange(r.Context(), db.RejectOpeningStateChangeReq{
				ID:         rejectReq.ID,
				EmployerID: orgUser.EmployerID,
				RejectedBy: orgUser.ID,
			})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", rejectReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNoStateChangeWaiting) {
				h.Dbg("no state change waiting", "id", rejectReq.ID)
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			h.Dbg("failed to reject state change", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("rejected state change", "id", rejectReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
