package interview

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
)

func EmployerRSVPInterview(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerRSVPInterview")
		var rsvpReq common.RSVPInterviewRequest
		if err := json.NewDecoder(r.Body).Decode(&rsvpReq); err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &rsvpReq) {
			h.Dbg("Validation failed", "rsvpReq", rsvpReq)
			return
		}
		h.Dbg("Validated", "rsvpReq", rsvpReq)

		err := h.DB().EmployerRSVPInterview(r.Context(), rsvpReq)
		if err != nil {
			h.Dbg("Error", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
