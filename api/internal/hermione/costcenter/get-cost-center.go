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

func GetCostCenter(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCostCenter")
		var getCostCenterReq vetchi.GetCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&getCostCenterReq)
		if err != nil {
			h.Dbg("failed to decode get cost center request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getCostCenterReq) {
			h.Dbg("validation failed", "getCostCenterReq", getCostCenterReq)
			return
		}
		h.Dbg("validated", "getCostCenterReq", getCostCenterReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		cc, err := h.DB().GetCCByName(r.Context(), db.GetCCByNameReq{
			Name:       getCostCenterReq.Name,
			EmployerID: orgUser.EmployerID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("CC not found", "name", getCostCenterReq.Name)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(cc)
		if err != nil {
			h.Err("failed to encode cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
