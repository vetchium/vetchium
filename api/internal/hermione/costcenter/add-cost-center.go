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
		var addCostCenterReq vetchi.AddCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&addCostCenterReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("AddCostCenterReq", "req", addCostCenterReq)

		if !h.Vator().Struct(w, &addCostCenterReq) {
			return
		}
		h.Dbg("AddCostCenterReq is valid")

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
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Added CostCenter", "CC", ccReq, "ID", costCenterID)

		w.WriteHeader(http.StatusOK)
	}
}
