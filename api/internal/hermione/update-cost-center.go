package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) updateCostCenter(w http.ResponseWriter, r *http.Request) {
	var updateCCRequest vetchi.UpdateCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&updateCCRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, updateCCRequest) {
		return
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	updateCCReq := db.UpdateCCReq{
		Name:       updateCCRequest.Name,
		Notes:      updateCCRequest.Notes,
		EmployerID: orgUser.EmployerID,
	}

	err := h.db.UpdateCostCenter(r.Context(), updateCCReq)
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
