package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func MyProfessionalEmails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement the business logic for getting professional emails
		// This would typically involve:
		// 1. Getting the user ID from the session
		// 2. Fetching all professional emails for this user from the database
		// 3. Converting them to the response format

		emails := []hub.ProfessionalEmail{} // This should be populated from the database
		h.Dbg("fetching professional emails")

		if err := json.NewEncoder(w).Encode(emails); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
