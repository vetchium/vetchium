package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
)

func MyOfficialEmails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		emails, err := h.DB().GetMyOfficialEmails(r.Context())
		if err != nil {
			h.Dbg("failed to get my official emails", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("my official emails", "emails", emails)

		if err := json.NewEncoder(w).Encode(emails); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
