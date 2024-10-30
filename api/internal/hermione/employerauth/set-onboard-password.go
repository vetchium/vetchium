package employerauth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func SetOnboardPassword(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var setOnboardPasswordReq vetchi.SetOnboardPasswordRequest
		err := json.NewDecoder(r.Body).Decode(&setOnboardPasswordReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Log().Info(
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
			h.Log().Error("Failed to hash password", "error", err)
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
