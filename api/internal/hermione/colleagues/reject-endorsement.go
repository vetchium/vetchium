package colleagues

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func RejectEndorsement(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RejectEndorsement")
		var rejectReq hub.RejectEndorsementRequest
		if err := json.NewDecoder(r.Body).Decode(&rejectReq); err != nil {
			h.Dbg("Error decoding request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &rejectReq) {
			h.Dbg("Invalid request", "rejectReq", rejectReq)
			return
		}
		h.Dbg("Validated", "rejectReq", rejectReq)

		err := h.DB().RejectEndorsement(r.Context(), rejectReq)
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("Application not found or not allowed", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("Error rejecting endorsement", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Rejected endorsement", "rejectReq", rejectReq)
		w.WriteHeader(http.StatusOK)
	}
}
