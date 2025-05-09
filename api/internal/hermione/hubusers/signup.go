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
				"link": "https://vetchium.com/hub/signup?token=" + token,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(req.Email)},
			// TODO: The subject should move to hedwig. Changing this should
			// also change the subject in 0027-hubsignup_test.go
			Subject: "Vetchium user signup invite",
		})
		if err != nil {
			h.Dbg("failed to generate invite email", "error", err)
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
			// Log the actual error for debugging purposes before specific checks
			h.Dbg(
				"Database call to SignupHubUser failed",
				"email",
				req.Email,
				"error",
				err,
			)

			if errors.Is(err, db.ErrDomainNotApprovedForSignup) {
				http.Error(
					w,
					"The domain of the provided email address is not approved for signup.",
					http.StatusUnprocessableEntity,
				) // 422
				return
			}
			if errors.Is(err, db.ErrInviteNotNeeded) {
				http.Error(
					w,
					"An account with this email may already exist, or an invite has already been sent.",
					http.StatusConflict,
				) // 409
				return
			}

			// Generic internal server error for other unhandled DB errors
			h.Err(
				"Unhandled database error during hub user signup",
				"email",
				req.Email,
				"original_error",
				err.Error(),
			)
			http.Error(
				w,
				"An unexpected error occurred while processing your signup request.",
				http.StatusInternalServerError,
			)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
