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

func DefunctCostCenter(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var defunctCostCenterRequest vetchi.DefunctCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&defunctCostCenterRequest)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &defunctCostCenterRequest) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		err = h.DB().DefunctCostCenter(
			r.Context(),
			db.DefunctReq{
				EmployerID: orgUser.EmployerID,
				Name:       defunctCostCenterRequest.Name,
			},
		)
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
