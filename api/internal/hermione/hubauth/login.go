package hubauth

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered hub login")

		// Simulate a random delay to avoid timing attacks
		<-time.After(time.Duration(rand.Intn(1000)) * time.Millisecond)

		var loginRequest vetchi.LoginRequest
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
		h.Dbg("validated request", "loginRequest", loginRequest)

		hubUser, err := h.DB().
			GetHubUserByEmail(r.Context(), string(loginRequest.Email))
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Dbg("failed to get hub user", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		if hubUser.State != vetchi.ActiveHubUserState {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(hubUser.PasswordHash),
			[]byte(loginRequest.Password),
		)
		if err != nil {
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
			Subject: "Vetchi Two Factor Authentication",
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
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		loginResponse := vetchi.LoginResponse{
			Token: tfaTokenString,
		}
		if err := json.NewEncoder(w).Encode(loginResponse); err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
