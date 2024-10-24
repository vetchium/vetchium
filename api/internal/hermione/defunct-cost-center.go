package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) defunctCostCenter(w http.ResponseWriter, r *http.Request) {
	var defunctCostCenterRequest vetchi.DefunctCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&defunctCostCenterRequest); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, defunctCostCenterRequest) {
		return
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	err := h.db.DefunctCostCenter(
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
