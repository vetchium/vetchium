package candidacy

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
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

		// Get candidate and opening details for email
		candidateInfo, err := h.DB().
			GetCandidateInfo(r.Context(), offerToCandidateRequest.CandidacyID)
		if err != nil {
			h.Dbg("Error getting candidate info", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Generate email using hedwig
		emailReq := hedwig.GenerateEmailReq{
			TemplateName: hedwig.NotifyCandidateOffer,
			Args: map[string]string{
				"CandidateName": candidateInfo.CandidateName,
				"CompanyName":   candidateInfo.CompanyName,
				"OpeningTitle":  candidateInfo.OpeningTitle,
				"CandidacyURL": fmt.Sprintf(
					"http://localhost:3001/candidacy/%s",
					offerToCandidateRequest.CandidacyID,
				),
			},
			EmailFrom: "notifications@vetchi.com",
			EmailTo:   []string{candidateInfo.CandidateEmail},
			Subject:   "Congratulations! Offer from " + candidateInfo.CompanyName,
		}

		email, err := h.Hedwig().GenerateEmail(emailReq)
		if err != nil {
			h.Dbg("Error generating email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Create request with comment and email
		req := db.OfferToCandidateReq{
			CandidacyID: offerToCandidateRequest.CandidacyID,
			Comment:     "Offer extended to candidate",
			Email:       email,
		}

		err = h.DB().OfferToCandidate(r.Context(), req)
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
