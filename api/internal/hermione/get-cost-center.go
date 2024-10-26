package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) GetCostCenter(w http.ResponseWriter, r *http.Request) {
	var getCostCenterReq vetchi.GetCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&getCostCenterReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, &getCostCenterReq) {
		return
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	cc, err := h.db.GetCCByName(r.Context(), db.GetCCByNameReq{
		Name:       getCostCenterReq.Name,
		EmployerID: orgUser.EmployerID,
	})

	if err != nil {
		if errors.Is(err, db.ErrNoCostCenter) {
			http.Error(w, "", http.StatusNotFound)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(cc)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
