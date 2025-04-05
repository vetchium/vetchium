package hubauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func ForgotPassword(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ForgotPassword")
		defer func() {
			h.Dbg("Sleep for random duration to avoid timing attacks")
			<-time.After(
				time.Millisecond * time.Duration(
					rand.Intn(int(h.Config().TimingAttackDelay.Milliseconds())),
				),
			)
		}()

		var forgotPasswordRequest hub.ForgotPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&forgotPasswordRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &forgotPasswordRequest) {
			h.Dbg("failed to validate request")
			return
		}
		h.Dbg("validated", "forgotPasswordRequest", forgotPasswordRequest)

		hubUser, err := h.DB().
			GetHubUserByEmail(r.Context(), string(forgotPasswordRequest.Email))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				// Send OK irrespective of whether the user exists or not
				http.Error(w, "", http.StatusOK)
				return
			}

			h.Dbg("failed to get hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		passwordResetToken := util.RandomString(
			vetchi.PasswordResetTokenLenBytes,
		)
		h.Dbg("token", "passwordResetToken", passwordResetToken)

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.HubPasswordReset,
			Args: map[string]string{
				"link": fmt.Sprintf(
					"%s/reset-password?token=%s",
					h.Config().Hub.WebURL,
					passwordResetToken,
				),
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(forgotPasswordRequest.Email)},

			// TODO: The subject should be from Hedwig, based on the template
			// This subject is used in 0007-hubauth_test.go too. Any change
			// in either place should be synced.
			Subject: "Vetchium Password Reset",
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().
			InitHubUserPasswordReset(r.Context(), db.HubUserInitPasswordReset{
				HubTokenReq: db.HubTokenReq{
					Token:            passwordResetToken,
					TokenType:        db.HubUserResetPasswordToken,
					ValidityDuration: h.Config().Hub.PasswordResetTokLife,
					HubUserID:        hubUser.ID,
				},
				Email: email,
			})
		if err != nil {
			h.Dbg("failed to init password reset", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
