package candidacy

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func OfferToCandidate(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var offerToCandidateRequest employer.OfferToCandidateRequest
		err := json.NewDecoder(r.Body).Decode(&offerToCandidateRequest)
		if err != nil {
			h.Dbg("Error decoding offer to candidate request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &offerToCandidateRequest) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("Validated")

		err = h.DB().OfferToCandidate(r.Context(), offerToCandidateRequest)
		if err != nil {
			if errors.Is(err, db.ErrNoCandidacy) {
				h.Dbg("Candidacy not found")
				http.Error(w, "", http.StatusNotFound)
				return
			}
			h.Dbg("Error offering to candidate", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("Offered to candidate")
		w.WriteHeader(http.StatusOK)
	}
}
