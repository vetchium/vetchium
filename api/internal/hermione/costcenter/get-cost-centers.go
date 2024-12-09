package costcenter

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetCostCenters(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCostCenters")
		var getCostCentersRequest employer.GetCostCentersRequest
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

		costCenters, err := h.DB().
			GetCostCenters(r.Context(), getCostCentersRequest)
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
