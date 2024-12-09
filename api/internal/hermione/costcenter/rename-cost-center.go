package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func RenameCostCenter(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RenameCostCenter")
		var renameCostCenterReq employer.RenameCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&renameCostCenterReq)
		if err != nil {
			h.Dbg("failed to decode rename cost center request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &renameCostCenterReq) {
			h.Dbg("validation failed", "renameCCReq", renameCostCenterReq)
			return
		}

		err = h.DB().RenameCostCenter(r.Context(), renameCostCenterReq)
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("CC not found", "name", renameCostCenterReq.OldName)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrDupCostCenterName) {
				h.Dbg("CC name exists", "name", renameCostCenterReq.NewName)
				http.Error(w, "", http.StatusConflict)
				return
			}

			h.Dbg("failed to rename cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("renamed cost center", "renameCostCenterReq", renameCostCenterReq)
		w.WriteHeader(http.StatusOK)
	}
}
