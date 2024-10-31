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
		var renameCostCenterReq vetchi.RenameCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&renameCostCenterReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &renameCostCenterReq) {
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
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrDupCostCenterName) {
				http.Error(w, "", http.StatusConflict)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
