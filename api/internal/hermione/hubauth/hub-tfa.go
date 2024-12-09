package hubauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/hub"
)

func HubTFA(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered HubTFA")
		var hubTFARequest hub.HubTFARequest
		err := json.NewDecoder(r.Body).Decode(&hubTFARequest)
		if err != nil {
			h.Err("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &hubTFARequest) {
			h.Dbg("failed to validate request", "error", err)
			return
		}
		h.Dbg("validated request", "hubTFARequest", hubTFARequest)

		hubUser, err := h.DB().GetHubUserByTFACreds(
			r.Context(),
			hubTFARequest.TFAToken,
			hubTFARequest.TFACode,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("no hub user found", "error", err)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Err("failed to get hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// We got the hub user. Now we need to create a session token
		// for the user. We are not deleting the TFA Code and the
		// TFA Token because if the /hub/tfa response could not be
		// delivered then the user will not be able to retry. But
		// remember to keep the TFA Token lifetime short.
		sessionToken := util.RandomString(vetchi.SessionTokenLenBytes)
		validityDuration := h.Config().Hub.SessionTokLife
		tokenType := db.HubUserSessionToken

		if hubTFARequest.RememberMe {
			tokenType = db.HubUserLTSToken
			validityDuration = h.Config().Hub.LTSTokLife
			h.Dbg("remember me", "validityDuration", validityDuration)
		}

		tokenReq := db.HubTokenReq{
			Token:            sessionToken,
			TokenType:        tokenType,
			ValidityDuration: validityDuration,
			HubUserID:        hubUser.ID,
		}

		err = h.DB().CreateHubUserToken(r.Context(), tokenReq)
		if err != nil {
			h.Err("failed to create hub user token", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hub.HubTFAResponse{
			SessionToken: sessionToken,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
