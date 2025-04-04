package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func GetCostCenter(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetCostCenter")
		var getCostCenterReq employer.GetCostCenterRequest
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

		cc, err := h.DB().GetCCByName(r.Context(), getCostCenterReq)
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
