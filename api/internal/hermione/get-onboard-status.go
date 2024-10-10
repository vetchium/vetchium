package hermione

import (
	"encoding/json"
	"net/http"

	"vetchi.org/pkg/libvetchi"
)

func (h *Hermione) getOnboardStatus(w http.ResponseWriter, r *http.Request) {
	var req libvetchi.GetOnboardStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	employer, err := h.db.GetEmployer(r.Context(), req.ClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := libvetchi.GetOnboardStatusResponse{
		Status: libvetchi.OnboardStatus(employer.OnboardStatus),
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		h.logger.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
