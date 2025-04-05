package hubauth

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
	"golang.org/x/crypto/bcrypt"
)

func Login(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered hub login")

		// Simulate a random delay to avoid timing attacks
		<-time.After(
			time.Millisecond * time.Duration(
				rand.Intn(int(h.Config().TimingAttackDelay.Milliseconds())),
			),
		)

		var loginRequest hub.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &loginRequest) {
			h.Dbg("failed to validate request", "error", err)
			return
		}
		h.Dbg("validated request")

		hubUser, err := h.DB().
			GetHubUserByEmail(r.Context(), string(loginRequest.Email))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("no hub user found", "email", loginRequest.Email)
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Dbg("failed to get hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if hubUser.State != hub.ActiveHubUserState {
			h.Dbg("hub user is not active", "email", loginRequest.Email)
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(hubUser.PasswordHash),
			[]byte(loginRequest.Password),
		)
		if err != nil {
			h.Dbg("invalid password", "email", loginRequest.Email)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		tfaMailCode, err := util.RandNumString(6)
		if err != nil {
			h.Dbg("failed to generate tfa mail code", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.HubUserTFA,
			Args:         map[string]string{"code": tfaMailCode},
			EmailFrom:    vetchi.EmailFrom,
			EmailTo:      []string{string(hubUser.Email)},

			// TODO: This subject should be from Hedwig, based on the template
			// This subject is used in 0007-hub-login_test.go too. Any change
			// in either place should be synced.
			Subject: "Vetchium Two Factor Authentication",
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		tfaTokenString := util.RandomString(vetchi.TGTokenLenBytes)

		err = h.DB().InitHubUserTFA(
			r.Context(),
			db.HubUserTFA{
				TFAToken: db.HubTokenReq{
					Token:            tfaTokenString,
					TokenType:        db.HubUserTFAToken,
					ValidityDuration: h.Config().Hub.TFATokLife,
					HubUserID:        hubUser.ID,
				},
				TFACode: tfaMailCode,
				Email:   email,
			},
		)
		if err != nil {
			h.Dbg("failed to init hub user tfa", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		loginResponse := hub.LoginResponse{
			Token: tfaTokenString,
		}
		if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
