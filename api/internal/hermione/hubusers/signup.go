package hubusers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func SignupHubUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered SignupHubUser")
		var req hub.SignupHubUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "request", req)
			return
		}
		h.Dbg("validated", "signupHubUserRequest", req)

		inviteEmail, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			// TODO: Fill the fields
		})
		if err != nil {
			h.Dbg("failed to generate invite email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		token := util.RandomUniqueID(vetchi.HubUserInviteTokenLenBytes)
		tokenValidTill := time.Now().Add(vetchi.HubUserInviteTokenValidity)

		err = h.DB().SignupHubUser(r.Context(), db.SignupHubUserReq{
			EmailAddress:   string(req.Email),
			InviteMail:     inviteEmail,
			Token:          token,
			TokenValidTill: tokenValidTill,
		})
		if err != nil {
			if errors.Is(err, db.ErrUnsupportedDomain) {
				h.Dbg("email domain is not supported for signup")
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to set handle", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
