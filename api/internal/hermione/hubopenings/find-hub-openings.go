package hubopenings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func FindHubOpenings(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FindHubOpenings")
		var findHubOpeningsReq hub.FindHubOpeningsRequest
		err := json.NewDecoder(r.Body).Decode(&findHubOpeningsReq)
		if err != nil {
			h.Dbg("failed to decode find hub openings request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &findHubOpeningsReq) {
			h.Dbg("validation failed", "findHubOpeningsReq", findHubOpeningsReq)
			return
		}
		h.Dbg("validated", "findHubOpeningsReq", findHubOpeningsReq)

		if findHubOpeningsReq.Limit == 0 {
			findHubOpeningsReq.Limit = 40
		}

		openings, err := h.DB().
			FindHubOpenings(r.Context(), &findHubOpeningsReq)
		if err != nil {
			h.Dbg("failed to find hub openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("found hub openings", "openings", openings)
		err = json.NewEncoder(w).Encode(openings)
		if err != nil {
			h.Err("failed to encode hub openings", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
