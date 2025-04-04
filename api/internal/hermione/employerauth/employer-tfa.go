package employerauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/employer"
)

func EmployerTFA(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered employer tfa")
		var employerTFARequest employer.EmployerTFARequest
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

		orgUser, err := h.DB().GetOrgUserByTFACreds(
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

		// We got the org user. Now we need to create a session token
		// for the user. We are not deleting the TFA Code and the
		// TFA Token because if the /employer/tfa response could not
		// be delivered then the user will not be able to retry. But
		// remember to keep the TFA Token lifetime short.
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

		err = json.NewEncoder(w).Encode(employer.EmployerTFAResponse{
			SessionToken: sessionToken,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
