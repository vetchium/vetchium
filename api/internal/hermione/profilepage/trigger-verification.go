package profilepage

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/hedwig"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/hub"
)

func TriggerVerification(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.TriggerVerificationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "req", req)
			return
		}
		h.Dbg("validated", "req", req)

		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Err("failed to get hub user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Check if the email exists and if it needs verification
		officialEmail, err := h.DB().
			GetOfficialEmail(r.Context(), string(req.Email))
		if err != nil {
			if err == db.ErrOfficialEmailNotFound {
				h.Dbg("email not found", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
			h.Dbg("failed to get official email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Check if last verification was more than 90 days ago
		if officialEmail.LastVerifiedAt != nil {
			timeSinceLastVerification := time.Now().
				UTC().
				Sub(*officialEmail.LastVerifiedAt)
			if timeSinceLastVerification < vetchi.VerificationValidityDuration {
				h.Dbg(
					"email was verified recently",
					"time_since_verification",
					timeSinceLastVerification,
				)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
		}

		code := util.RandomString(vetchi.AddOfficialEmailCodeLenBytes)

		_, err = h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.AddOfficialEmail,
			Args: map[string]string{
				"Name":   hubUser.FullName,
				"Handle": hubUser.Handle,
				"Code":   code,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(req.Email)},
			Subject:   "Vetchium - Confirm Email Ownership",
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().
			UpdateOfficialEmailVerificationCode(r.Context(), db.UpdateOfficialEmailVerificationCodeReq{
				Email:   string(req.Email),
				Code:    code,
				HubUser: hubUser,
			})
		if err != nil {
			if err == db.ErrOfficialEmailNotFound {
				h.Dbg("email was deleted before update", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}
			h.Dbg("failed to update verification code", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("triggered verification", "email", req.Email)
		w.WriteHeader(http.StatusOK)
	}
}
