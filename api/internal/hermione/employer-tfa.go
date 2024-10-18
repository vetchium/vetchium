package hermione

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) employerTFA(w http.ResponseWriter, r *http.Request) {
	var employerTFARequest vetchi.EmployerTFARequest

	err := json.NewDecoder(r.Body).Decode(&employerTFARequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, employerTFARequest) {
		return
	}

	// TODO: Validate token

}
