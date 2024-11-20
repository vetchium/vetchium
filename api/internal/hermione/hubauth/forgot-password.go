package hubauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ForgotPasswordHandler(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			h.Dbg("Sleep for random duration to avoid timing attacks")
			<-time.After(
				time.Millisecond * time.Duration(
					rand.Intn(int(h.Config().TimingAttackDelay.Milliseconds())),
				),
			)
		}()

		h.Dbg("Entered ForgotPasswordHandler")

		var forgotPasswordRequest vetchi.ForgotPasswordRequest
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
			Subject:   "Reset your Vetchi password",
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().
			InitHubUserPasswordReset(r.Context(), db.HubUserPasswordResetReq{
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
