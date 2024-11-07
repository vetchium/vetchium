package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ApproveOpeningStateChange(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ApproveOpeningStateChange")
		var approveReq vetchi.ApproveOpeningStateChangeRequest
		err := json.NewDecoder(r.Body).Decode(&approveReq)
		if err != nil {
			h.Dbg("failed to decode approve state change request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &approveReq) {
			h.Dbg("validation failed", "approveReq", approveReq)
			return
		}
		h.Dbg("validated", "approveReq", approveReq)

		err = h.DB().ApproveOpeningStateChange(r.Context(), approveReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", approveReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrNoStateChangeWaiting) {
				h.Dbg("no state change waiting", "id", approveReq.ID)
				http.Error(w, "", http.StatusBadRequest)
				return
			}

			h.Dbg("failed to approve state change", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("approved state change", "id", approveReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
