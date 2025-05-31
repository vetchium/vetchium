package employerauth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"

	"golang.org/x/crypto/bcrypt"
)

func ResetPassword(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ResetPassword")
		var resetPasswordRequest employer.EmployerResetPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&resetPasswordRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("reset password request", "request", resetPasswordRequest)

		if !h.Vator().Struct(w, &resetPasswordRequest) {
			h.Dbg("invalid request", "request", resetPasswordRequest)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("validated", "resetPasswordRequest", resetPasswordRequest)

		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(resetPasswordRequest.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			h.Dbg("failed to hash password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		employerPasswordReset := db.EmployerPasswordReset{
			Token:        resetPasswordRequest.Token,
			PasswordHash: string(passwordHash),
		}

		err = h.DB().ResetEmployerPassword(r.Context(), employerPasswordReset)
		if err != nil {
			if errors.Is(err, db.ErrInvalidPasswordResetToken) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			h.Dbg("failed to reset password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
