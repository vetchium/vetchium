package orgusers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"

	"golang.org/x/crypto/bcrypt"
)

func SignupOrgUser(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered SignupOrgUser")
		var signupOrgUserRequest employer.SignupOrgUserRequest
		err := json.NewDecoder(r.Body).Decode(&signupOrgUserRequest)
		if err != nil {
			h.Dbg("failed to decode signup org user request", "err", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("SignupOrgUserRequest", "req", signupOrgUserRequest)

		if !h.Vator().Struct(w, &signupOrgUserRequest) {
			h.Dbg("validation failed", "req", signupOrgUserRequest)
			return
		}
		h.Dbg("validated", "SignupOrgUserRequest", signupOrgUserRequest)

		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(signupOrgUserRequest.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			h.Dbg("failed to hash password", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().SignupOrgUser(r.Context(), db.SignupOrgUserReq{
			InviteToken:  signupOrgUserRequest.InviteToken,
			Name:         signupOrgUserRequest.Name,
			PasswordHash: string(passwordHash),
		})
		if err != nil {
			if errors.Is(err, db.ErrInviteTokenNotFound) {
				h.Dbg("Not Found", "token", signupOrgUserRequest.InviteToken)
				http.Error(w, "", http.StatusForbidden)
				return
			}

			h.Dbg("failed to signup org user", "err", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
