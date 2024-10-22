package hermione

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) getCostCenters(w http.ResponseWriter, r *http.Request) {
	var getCostCentersRequest vetchi.GetCostCentersRequest
	if err := json.NewDecoder(r.Body).Decode(&getCostCentersRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, getCostCentersRequest) {
		return
	}

	employerIDString := r.Header.Get(middleware.EmployerIDHeader)
	if employerIDString == "" {
		h.log.Error("X-Vetchi-EmployerID header missing")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	employerID, err := uuid.Parse(employerIDString)
	if err != nil {
		h.log.Error("Invalid X-Vetchi-EmployerID", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	costCenters, err := h.db.GetCostCenters(
		r.Context(),
		db.CCentersList{
			EmployerID: employerID,
			Offset:     getCostCentersRequest.Offset,
			Limit:      getCostCentersRequest.Limit,
		},
	)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(costCenters)
	if err != nil {
		h.log.Error("Error encoding response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
