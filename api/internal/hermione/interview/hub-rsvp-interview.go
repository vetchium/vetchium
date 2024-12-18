package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func HubRSVPInterview(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered HubRSVPInterview")
		var rsvpReq hub.HubRSVPInterviewRequest
		if err := json.NewDecoder(r.Body).Decode(&rsvpReq); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &rsvpReq) {
			return
		}
		h.Dbg("Validated", "rsvpReq", rsvpReq)

		err := h.DB().HubRSVPInterview(r.Context(), rsvpReq)
		if err != nil {
			if errors.Is(err, db.ErrNoInterview) {
				http.Error(w, "", http.StatusNotFound)
				return
			}
			if errors.Is(err, db.ErrInvalidInterviewState) {
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
			h.Dbg("RSVP status update failed", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("RSVP status updated")

		w.WriteHeader(http.StatusOK)
	}
}
