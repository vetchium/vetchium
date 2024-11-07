package openings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func FilterOpenings(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterOpenings")
		var filterOpeningsReq vetchi.FilterOpeningsRequest
		err := json.NewDecoder(r.Body).Decode(&filterOpeningsReq)
		if err != nil {
			h.Dbg("failed to decode filter openings request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterOpeningsReq) {
			h.Dbg("validation failed", "filterOpeningsReq", filterOpeningsReq)
			return
		}
		h.Dbg("validated", "filterOpeningsReq", filterOpeningsReq)

		if filterOpeningsReq.Limit == 0 {
			filterOpeningsReq.Limit = 40
			h.Dbg("set default limit", "limit", filterOpeningsReq.Limit)
		}

		if len(filterOpeningsReq.State) == 0 {
			filterOpeningsReq.State = []vetchi.OpeningState{
				vetchi.ActiveOpening,
			}
			h.Dbg("set default state", "state", filterOpeningsReq.State)
		}

		states := make([]string, len(filterOpeningsReq.State))
		for i, state := range filterOpeningsReq.State {
			states[i] = string(state)
		}

		openings, err := h.DB().FilterOpenings(r.Context(), filterOpeningsReq)
		if err != nil {
			h.Dbg("failed to filter openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(openings)
		if err != nil {
			h.Err("failed to encode openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
