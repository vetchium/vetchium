package interview

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func GetHubInterviewsByCandidacy(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetHubInterviewsByCandidacy")
		var req hub.GetHubInterviewsByCandidacyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("request decoded", "request", req)

		interviews, err := h.DB().GetHubInterviewsByCandidacy(r.Context(), req)
		if err != nil {
			h.Dbg("failed to get hub interviews by candidacy", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		h.Dbg("hub interviews by candidacy", "interviews", interviews)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(interviews)
	}
}
