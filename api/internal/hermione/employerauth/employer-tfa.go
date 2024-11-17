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
		h.Dbg("Entered employer tfa")
		var employerTFARequest vetchi.EmployerTFARequest
		err := json.NewDecoder(r.Body).Decode(&employerTFARequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &employerTFARequest) {
			h.Dbg("failed to validate request", "error", err)
			return
		}
		h.Dbg("validated request", "employerTFARequest", employerTFARequest)

		orgUser, err := h.DB().GetOrgUserByToken(
			r.Context(),
			employerTFARequest.TFACode,
			employerTFARequest.TFAToken,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				h.Dbg("no org user found", "error", err)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Dbg("failed to get org user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		sessionToken := util.RandomString(vetchi.SessionTokenLenBytes)
		validityDuration := h.Config().Employer.SessionTokLife
		tokenType := db.EmployerSessionToken

		if employerTFARequest.RememberMe {
			tokenType = db.EmployerLTSToken
			validityDuration = h.Config().Employer.LTSTokLife
			h.Dbg("remember me", "validityDuration", validityDuration)
		}

		tokenReq := db.EmployerTokenReq{
			Token:            sessionToken,
			TokenType:        tokenType,
			ValidityDuration: validityDuration,
			OrgUserID:        orgUser.ID,
		}
		h.Dbg("creating org user token", "tokenReq", tokenReq)

		err = h.DB().CreateOrgUserToken(r.Context(), tokenReq)
		if err != nil {
			h.Dbg("failed to create org user token", "error", err)
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
