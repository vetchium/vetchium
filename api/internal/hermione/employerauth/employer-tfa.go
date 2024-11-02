package employerauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func EmployerTFA(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var employerTFARequest vetchi.EmployerTFARequest

		err := json.NewDecoder(r.Body).Decode(&employerTFARequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &employerTFARequest) {
			return
		}

		orgUser, err := h.DB().GetOrgUserByToken(
			r.Context(),
			employerTFARequest.TFACode,
			employerTFARequest.TFAToken,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		sessionToken := util.RandomString(vetchi.SessionTokenLenBytes)
		validityDuration, err := h.ConfigDuration(db.EmployerSessionToken)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		tokenType := db.EmployerSessionToken

		if employerTFARequest.RememberMe {
			tokenType = db.EmployerLTSToken
			validityDuration, err = h.ConfigDuration(db.EmployerLTSToken)
			if err != nil {
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
		}

		err = h.DB().CreateOrgUserToken(r.Context(), db.TokenReq{
			Token:            sessionToken,
			TokenType:        tokenType,
			ValidityDuration: validityDuration,
			OrgUserID:        orgUser.ID,
		})
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(vetchi.EmployerTFAResponse{
			SessionToken: sessionToken,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
