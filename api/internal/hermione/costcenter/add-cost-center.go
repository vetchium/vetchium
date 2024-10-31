package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddCostCenter(h vhandler.VHandler) http.HandlerFunc {
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		ccReq := db.CCenterReq{
			Name:       addCostCenterReq.Name,
			Notes:      addCostCenterReq.Notes,
			EmployerID: orgUser.EmployerID,
			OrgUserID:  orgUser.ID,
		}

		costCenterID, err := h.DB().CreateCostCenter(r.Context(), ccReq)
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

		h.Dbg("Added CostCenter", "CC", ccReq, "ID", costCenterID)
		w.WriteHeader(http.StatusOK)
	}
}
