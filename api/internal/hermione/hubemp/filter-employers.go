package hubemp

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func FilterEmployers(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterEmployers")

		var filterEmployersReq hub.FilterEmployersRequest
		err := json.NewDecoder(r.Body).Decode(&filterEmployersReq)
		if err != nil {
			h.Dbg("failed to decode filter employers request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		h.Dbg("FilterEmployersRequest", "request", filterEmployersReq)

		employers, err := h.DB().
			FilterEmployers(r.Context(), filterEmployersReq)
		if err != nil {
			h.Dbg("failed to get employers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Filtered Employers", "employers", employers)

		err = json.NewEncoder(w).Encode(hub.FilterEmployersResponse{
			Employers: employers,
		})
		if err != nil {
			h.Dbg("failed to encode employers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
