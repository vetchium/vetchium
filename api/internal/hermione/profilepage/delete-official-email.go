package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func DeleteOfficialEmail(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.DeleteOfficialEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "req", req)
			return
		}
		h.Dbg("validated", "req", req)

		// TODO: Implement the business logic for deleting official email
		// This would typically involve:
		// 1. Validating that the email exists
		// 2. Checking if the user has permission to delete this email
		// 3. Removing the email from the database

		w.WriteHeader(http.StatusOK)
	}
}
