package hermione

import (
	"encoding/json"
	"net/http"

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

	if getCostCentersRequest.Limit <= 0 {
		getCostCentersRequest.Limit = 100
	}

	if len(getCostCentersRequest.States) == 0 {
		getCostCentersRequest.States = []vetchi.CostCenterState{vetchi.ActiveCC}
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	states := []string{}
	for _, state := range getCostCentersRequest.States {
		states = append(states, string(state))
	}

	costCenters, err := h.db.GetCostCenters(
		r.Context(),
		db.CCentersList{
			EmployerID: orgUser.EmployerID,
			Offset:     getCostCentersRequest.Offset,
			Limit:      getCostCentersRequest.Limit,
			States:     states,
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
