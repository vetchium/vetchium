package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func MyCandidacies(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered MyCandidacies")

		var getMyCandidaciesReq hub.MyCandidaciesRequest
		if err := json.NewDecoder(r.Body).Decode(&getMyCandidaciesReq); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getMyCandidaciesReq) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated request", "getMyCandidaciesReq", getMyCandidaciesReq)

		if getMyCandidaciesReq.Limit == 0 {
			getMyCandidaciesReq.Limit = 40
		}

		candidacies, err := h.DB().
			GetMyCandidacies(r.Context(), getMyCandidaciesReq)
		if err != nil {
			h.Dbg("Error getting candidacies", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got candidacies", "candidacies", candidacies)

		err = json.NewEncoder(w).Encode(candidacies)
		if err != nil {
			h.Err("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
