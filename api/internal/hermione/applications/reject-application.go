package applications

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func RejectApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RejectApplication")
		var rejectApplicationReq vetchi.RejectApplicationRequest
		err := json.NewDecoder(r.Body).Decode(&rejectApplicationReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &rejectApplicationReq) {
			h.Dbg("failed to validate request", "error", err)
			return
		}
		h.Dbg("validated", "rejectApplicationReq", rejectApplicationReq)

		mailInfo, err := h.DB().
			GetApplicationMailInfo(r.Context(), rejectApplicationReq.ApplicationID)
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("application not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}
			h.Dbg("failed to get application mail info", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.RejectApplication,
			Args: map[string]string{
				"hub_user_full_name":      mailInfo.HubUser.FullName,
				"employer_company_name":   mailInfo.Employer.CompanyName,
				"employer_primary_domain": mailInfo.Employer.PrimaryDomain,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{mailInfo.HubUser.Email},
			Subject: fmt.Sprintf(
				"%s - Rejected",
				mailInfo.Employer.CompanyName,
			),
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().RejectApplication(r.Context(), db.RejectApplicationRequest{
			ApplicationID: rejectApplicationReq.ApplicationID,
			Email:         email,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("application not found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrApplicationStateInCompatible) {
				h.Dbg("application is not in applied state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to reject application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
