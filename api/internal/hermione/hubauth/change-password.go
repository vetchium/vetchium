package hubauth

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangePassword")

		var changePasswordRequest hub.ChangePasswordRequest
		err := json.NewDecoder(r.Body).Decode(&changePasswordRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &changePasswordRequest) {
			h.Dbg("failed to validate request")
			return
		}

		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Err("failed to get hub user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Verify old password
		err = bcrypt.CompareHashAndPassword(
			[]byte(hubUser.PasswordHash),
			[]byte(changePasswordRequest.OldPassword),
		)
		if err != nil {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		// Hash new password
		newPasswordHash, err := bcrypt.GenerateFromPassword(
			[]byte(changePasswordRequest.NewPassword),
			bcrypt.DefaultCost,
		)
		if err != nil {
			h.Err("failed to hash password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().
			ChangeHubUserPassword(r.Context(), hubUser.ID, string(newPasswordHash))
		if err != nil {
			h.Err("failed to change password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
