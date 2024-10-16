package hermione

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"golang.org/x/crypto/bcrypt"
)

func (h *Hermione) setOnboardPassword(w http.ResponseWriter, r *http.Request) {
	var req vetchi.SetOnboardPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Validate the password

	// Hash the password using bcrypt
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	err = h.db.OnboardAdmin(
		r.Context(),
		req.ClientID,
		string(passwordHash),
		req.Token,
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
