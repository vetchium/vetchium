package employerauth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"

	"golang.org/x/crypto/bcrypt"
)

func SetOnboardPassword(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var setOnboardPasswordReq employer.SetOnboardPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&setOnboardPasswordReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Inf(
			"Set Onboard Password Request",
			"request",
			setOnboardPasswordReq,
		)
		log.Printf("Set Onboard Password Request %+v", setOnboardPasswordReq)

		if !h.Vator().Struct(w, &setOnboardPasswordReq) {
			return
		}

		// Hash the password using bcrypt
		passwordHash, err := bcrypt.GenerateFromPassword(
			[]byte(setOnboardPasswordReq.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			h.Err("Failed to hash password", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().OnboardAdmin(
			r.Context(),
			db.OnboardReq{
				DomainName: setOnboardPasswordReq.ClientID,
				Password:   string(passwordHash),
				Token:      setOnboardPasswordReq.Token,
			},
		)
		if err != nil {
			if err == db.ErrNoEmployer || err == db.ErrOrgUserAlreadyExists {
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
