package employerauth

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/middleware"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
	"golang.org/x/crypto/bcrypt"
)

func ChangePassword(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ChangePassword")

		var changePasswordRequest employer.EmployerChangePasswordRequest
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get org user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// Verify old password
		err = bcrypt.CompareHashAndPassword(
			[]byte(orgUser.PasswordHash),
			[]byte(changePasswordRequest.OldPassword),
		)
		if err != nil {
			h.Dbg("failed to verify old password", "error", err)
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
			ChangeOrgUserPassword(r.Context(), orgUser.ID, string(newPasswordHash))
		if err != nil {
			h.Err("failed to change password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("password changed successfully")
		w.WriteHeader(http.StatusOK)
	}
}
