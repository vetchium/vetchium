package applications

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ShortlistApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ShortlistApplication")
		var shortlistRequest vetchi.ShortlistApplicationRequest
		err := json.NewDecoder(r.Body).Decode(&shortlistRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("shortlistRequest", "request", shortlistRequest)

		if !h.Vator().Struct(w, &shortlistRequest) {
			h.Dbg("invalid request", "error", err)
			return
		}
		h.Dbg("validated", "shortlistRequest", shortlistRequest)

		/*
			Get the primary domain and company name for the given application_id
			Generate an email to notify the candidate, using the primary domain and company name
			Generate a CandidacyID
			In a db transaction,
				Create a Candidacy for the shortlisted candidate
				Update the application state to Shortlisted
		*/

		// Ensures secrecy
		candidacyID := util.RandomString(vetchi.CandidacyIDLenBytes)
		// Ensures uniqueness
		candidacyID = candidacyID + strconv.FormatInt(
			time.Now().UnixNano(),
			36,
		)

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: "shortlist_application",
		})

		err = h.DB().ShortlistApplication(r.Context(), db.ShortlistRequest{
			ApplicationID: shortlistRequest.ApplicationID,
			CandidacyID:   candidacyID,
			Email:         email,
		})
		if err != nil {
			h.Dbg("failed to shortlist application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
