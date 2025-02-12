package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func MyOfficialEmails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: Implement the business logic for getting official emails
		// This would typically involve:
		// 1. Getting the user ID from the session
		// 2. Fetching all official emails for this user from the database
		// 3. Converting them to the response format

		emails := []hub.OfficialEmail{} // This should be populated from the database
		h.Dbg("fetching official emails")

		if err := json.NewEncoder(w).Encode(emails); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
