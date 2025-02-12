package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func TriggerVerification(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.TriggerVerificationRequest
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

		// TODO: Implement the business logic for triggering verification
		// This would typically involve:
		// 1. Validating that the email exists in our system
		// 2. Generating a verification code
		// 3. Sending the verification email
		// 4. Updating the VerifyInProgress status

		w.WriteHeader(http.StatusOK)
	}
}
