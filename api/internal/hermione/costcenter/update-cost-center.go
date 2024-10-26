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

func UpdateCostCenter(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateCCRequest vetchi.UpdateCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&updateCCRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateCCRequest) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		updateCCReq := db.UpdateCCReq{
			Name:       updateCCRequest.Name,
			Notes:      updateCCRequest.Notes,
			EmployerID: orgUser.EmployerID,
		}

		err = h.DB().UpdateCostCenter(r.Context(), updateCCReq)
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else if errors.Is(err, db.ErrDupCostCenterName) {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
