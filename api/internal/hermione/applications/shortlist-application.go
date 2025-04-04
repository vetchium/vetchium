package applications

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/employer"
)

func ShortlistApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ShortlistApplication")
		var shortlistRequest employer.ShortlistApplicationRequest
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

		mailInfo, err := h.DB().
			GetApplicationMailInfo(r.Context(), shortlistRequest.ApplicationID)
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("not found", "id", shortlistRequest.ApplicationID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get application mail info", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Ensures secrecy
		candidacyID := util.RandomString(vetchi.CandidacyIDLenBytes)
		// Ensures uniqueness
		candidacyID = candidacyID + strconv.FormatInt(
			time.Now().UnixNano(),
			36,
		)
		h.Dbg("New candidacyID generated", "id", candidacyID)

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.ShortlistApplication,
			Args: map[string]string{
				"hub_user_full_name":      mailInfo.HubUser.FullName,
				"employer_company_name":   mailInfo.Employer.CompanyName,
				"employer_primary_domain": mailInfo.Employer.PrimaryDomain,
				"candidacy_link":          vetchi.HubBaseURL + "/candidacy/" + candidacyID,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{mailInfo.HubUser.Email},
			Subject: fmt.Sprintf(
				"Shortlisted for %s",
				mailInfo.Employer.CompanyName,
			),
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().ShortlistApplication(r.Context(), db.ShortlistRequest{
			ApplicationID: shortlistRequest.ApplicationID,
			OpeningID:     mailInfo.Opening.OpeningID,
			CandidacyID:   candidacyID,
			Email:         email,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("failed to shortlist application", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrApplicationStateInCompatible) {
				h.Dbg("failed to shortlist application", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to shortlist application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
