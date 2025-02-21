package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func MyColleagueApprovals(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered MyColleagueApprovals")

		var req hub.MyColleagueApprovalsRequest
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

		// TODO: Implement DB call to get approvals
		approvals := []string{} // Replace with actual data

		if err := json.NewEncoder(w).Encode(approvals); err != nil {
			h.Err("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
