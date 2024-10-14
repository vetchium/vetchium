package hermione

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/libvetchi"
)

func (h *Hermione) setOnboardPassword(w http.ResponseWriter, r *http.Request) {
	var req libvetchi.SetOnboardPasswordRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Validate the password

	err = h.db.OnboardAdmin(
		r.Context(),
		req.ClientID,
		req.Password,
		req.Token,
	)
	if err != nil {
		if err == db.ErrNoEmployer {
			http.Error(w, "", http.StatusUnprocessableEntity)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
