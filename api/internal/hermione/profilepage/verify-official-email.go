package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func VerifyOfficialEmail(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.VerifyOfficialEmailRequest
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

		// TODO: Implement the business logic for verifying official email
		// This would typically involve:
		// 1. Validating the verification code
		// 2. Updating the email verification status in the database
		// 3. Setting the LastVerifiedAt timestamp

		w.WriteHeader(http.StatusOK)
	}
}
