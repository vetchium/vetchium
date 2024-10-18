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

	// TODO: Validate incoming tokens and generate a session token
	var employerTFAResponse vetchi.EmployerTFAResponse
	employerTFAResponse.SessionToken = "TODO: Hardcoded session token"

	err = json.NewEncoder(w).Encode(employerTFAResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
