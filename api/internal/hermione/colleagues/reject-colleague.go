package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func RejectColleague(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RejectColleague")

		var req hub.RejectColleagueRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated request", "req", req)

		// TODO: Implement DB call
		w.WriteHeader(http.StatusOK)
	}
}
