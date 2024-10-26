package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (h *Hermione) renameCostCenter(w http.ResponseWriter, r *http.Request) {
	var renameCostCenterReq vetchi.RenameCostCenterRequest
	if err := json.NewDecoder(r.Body).Decode(&renameCostCenterReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !h.vator.Struct(w, &renameCostCenterReq) {
		return
	}

	orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
	if !ok {
		h.log.Error("failed to get orgUser from context")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	renameCCReq := db.RenameCCReq{
		OldName:    renameCostCenterReq.OldName,
		NewName:    renameCostCenterReq.NewName,
		EmployerID: orgUser.EmployerID,
		OrgUserID:  orgUser.ID,
	}
	err := h.db.RenameCostCenter(r.Context(), renameCCReq)
	if err != nil {
		if errors.Is(err, db.ErrDupCostCenterName) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	return
}
