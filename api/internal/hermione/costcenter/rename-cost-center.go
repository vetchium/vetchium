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

func RenameCostCenter(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RenameCostCenter")
		var renameCostCenterReq vetchi.RenameCostCenterRequest
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		renameCCReq := db.RenameCCReq{
			OldName:    renameCostCenterReq.OldName,
			NewName:    renameCostCenterReq.NewName,
			EmployerID: orgUser.EmployerID,
			OrgUserID:  orgUser.ID,
		}

		err = h.DB().RenameCostCenter(r.Context(), renameCCReq)
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

		h.Dbg("renamed cost center", "renameCCReq", renameCCReq)
		w.WriteHeader(http.StatusOK)
	}
}
