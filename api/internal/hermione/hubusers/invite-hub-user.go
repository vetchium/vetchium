package hubusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func InviteHubUser(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered InviteHubUser")
		var hubUserInviteReq hub.HubUserInviteRequest
		err := json.NewDecoder(r.Body).Decode(&hubUserInviteReq)
		if err != nil {
			h.Dbg("failed to unmarshal request", "err", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &hubUserInviteReq) {
			h.Dbg("validation failed", "hubUserInviteReq", hubUserInviteReq)
			return
		}
		h.Dbg("validated", "hubUserInviteReq", hubUserInviteReq)

		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Err("failed to get hub user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		token := util.RandomUniqueID(vetchi.HubUerInviteTokenLenBytes)

		inviteMail, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.InviteHubUser,
			Args: map[string]string{
				"inviter": hubUser.FullName,
				"link":    vetchi.SignupHubUserURL + token,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(hubUserInviteReq.Email)},
		})
		if err != nil {
			h.Dbg("failed to generate email", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().InviteHubUser(r.Context(), db.InviteHubUserReq{
			EmailAddress: hubUserInviteReq.Email,
			InviteMail:   inviteMail,
			Token:        token,
		})
		if err != nil {
			if errors.Is(err, db.ErrInviteNotNeeded) {
				// Either the user is already a hubuser
				// or an invite was sent recently
				h.Dbg("New invite not needed")
				w.WriteHeader(http.StatusOK)
				return
			}

			h.Dbg("failed to add hub user invite", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("hubuser invite mail added", "hubUserInviteReq", hubUserInviteReq)
		w.WriteHeader(http.StatusOK)
	})
}
