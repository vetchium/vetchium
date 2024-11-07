package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddCostCenter(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddCostCenter")
		var addCostCenterReq vetchi.AddCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&addCostCenterReq)
		if err != nil {
			h.Dbg("failed to decode add cost center request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addCostCenterReq) {
			h.Dbg("validation failed", "addCostCenterReq", addCostCenterReq)
			return
		}
		h.Dbg("validated", "addCostCenterReq", addCostCenterReq)

		costCenterID, err := h.DB().
			CreateCostCenter(r.Context(), addCostCenterReq)
		if err != nil {
			if errors.Is(err, db.ErrDupCostCenterName) {
				h.Dbg("cost center exists", "name", addCostCenterReq.Name)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("failed to create cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Added CostCenter", "ID", costCenterID)
		w.WriteHeader(http.StatusOK)
	}
}
