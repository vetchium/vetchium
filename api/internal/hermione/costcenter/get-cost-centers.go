package costcenter

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func GetCostCenters(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCostCenters")
		var getCostCentersRequest vetchi.GetCostCentersRequest
		err := json.NewDecoder(r.Body).Decode(&getCostCentersRequest)
		if err != nil {
			h.Dbg("failed to decode get cost centers request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getCostCentersRequest) {
			h.Dbg("validation failed", "getCCsReq", getCostCentersRequest)
			return
		}
		h.Dbg("validated", "getCCsReq", getCostCentersRequest)

		if getCostCentersRequest.Limit <= 0 {
			getCostCentersRequest.Limit = 100
			h.Dbg("set default limit", "limit", getCostCentersRequest.Limit)
		}

		if len(getCostCentersRequest.States) == 0 {
			getCostCentersRequest.States = []vetchi.CostCenterState{
				vetchi.ActiveCC,
			}
			h.Dbg("set default states", "states", getCostCentersRequest.States)
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		states := []string{}
		for _, state := range getCostCentersRequest.States {
			// already validated by vator
			states = append(states, string(state))
		}

		costCenters, err := h.DB().GetCostCenters(
			r.Context(),
			db.CCentersList{
				EmployerID:    orgUser.EmployerID,
				States:        states,
				PaginationKey: getCostCentersRequest.PaginationKey,
				Limit:         getCostCentersRequest.Limit,
			},
		)
		if err != nil {
			h.Dbg("failed to get cost centers", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(costCenters)
		if err != nil {
			h.Err("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
