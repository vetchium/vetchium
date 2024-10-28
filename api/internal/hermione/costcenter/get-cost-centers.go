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
		var getCostCentersRequest vetchi.GetCostCentersRequest
		err := json.NewDecoder(r.Body).Decode(&getCostCentersRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getCostCentersRequest) {
			return
		}

		if getCostCentersRequest.Limit <= 0 {
			getCostCentersRequest.Limit = 100
		}

		if len(getCostCentersRequest.States) == 0 {
			getCostCentersRequest.States = []vetchi.CostCenterState{
				vetchi.ActiveCC,
			}
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
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
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(costCenters)
		if err != nil {
			h.Log().Error("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
