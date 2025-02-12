package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func AddProfessionalEmail(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.AddProfessionalEmailRequest
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

		// TODO: Implement the business logic for adding professional email
		// This would typically involve:
		// 1. Validating the email format
		// 2. Checking if the email already exists
		// 3. Adding the email to the database
		// 4. Potentially triggering a verification email

		w.WriteHeader(http.StatusOK)
	}
}
