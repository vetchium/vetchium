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

		token := util.RandomUniqueID(vetchi.HubUserInviteTokenLenBytes)
		tokenValidTill := time.Now().Add(vetchi.HubUserInviteTokenValidity)

		inviteEmail, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.InviteHubUserSignup,
			Args: map[string]string{
				"link": h.Config().SignupHubUserURL + token,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(req.Email)},
			// TODO: The subject should move to hedwig. Changing this should
			// also change the subject in 0027-hubsignup_test.go
			Subject: "Vetchium user signup invite",
		})
		if err != nil {
			h.Err("failed to generate invite email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().SignupHubUser(r.Context(), db.SignupHubUserReq{
			EmailAddress:   string(req.Email),
			InviteMail:     inviteEmail,
			Token:          token,
			TokenValidTill: tokenValidTill,
		})
		if err != nil {
			if errors.Is(err, db.ErrDomainNotApprovedForSignup) {
				h.Dbg("domain not approved for signup")
				http.Error(w, "", 460)
				return
			}

			if errors.Is(err, db.ErrInviteNotNeeded) {
				h.Dbg("invite not needed")
				http.Error(w, "", 461)
				return
			}

			// Generic internal server error for other unhandled DB errors
			h.Dbg("database error", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
