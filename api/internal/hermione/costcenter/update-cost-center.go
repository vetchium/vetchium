package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func UpdateCostCenter(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateCostCenter")
		var updateCCRequest employer.UpdateCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&updateCCRequest)
		if err != nil {
			h.Dbg("failed to decode update cost center request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateCCRequest) {
			h.Dbg("validation failed", "updateCCReq", updateCCRequest)
			return
		}
		h.Dbg("validated", "updateCCReq", updateCCRequest)

		err = h.DB().UpdateCostCenter(r.Context(), updateCCRequest)
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("cc not found", "name", updateCCRequest.Name)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("updated cost center", "updateCCReq", updateCCRequest)
		w.WriteHeader(http.StatusOK)
	}
}
