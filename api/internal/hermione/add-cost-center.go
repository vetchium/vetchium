package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) addCostCenter(w http.ResponseWriter, r *http.Request) {
	var addCostCenterReq vetchi.AddCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&addCostCenterReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, &addCostCenterReq) {
		return
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	ccReq := db.CCenterReq{
		Name:       addCostCenterReq.Name,
		Notes:      addCostCenterReq.Notes,
		EmployerID: orgUser.EmployerID,
		OrgUserID:  orgUser.ID,
	}

	costCenterID, err := h.db.CreateCostCenter(r.Context(), ccReq)
	if err != nil {
		if errors.Is(err, db.ErrDupCostCenterName) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	h.log.Debug("Created CostCenter", "CC", ccReq, "ID", costCenterID)

	err = json.NewEncoder(w).Encode(vetchi.AddCostCenterResponse{
		Name: addCostCenterReq.Name,
	})
	if err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
