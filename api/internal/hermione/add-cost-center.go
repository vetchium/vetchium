package hermione

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
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

	if !h.vator.Struct(w, addCostCenterReq) {
		return
	}

	orgUserIDString := r.Header.Get(middleware.OrgUserIDHeader)
	if orgUserIDString == "" {
		h.log.Error("X-Vetchi-OrgUserID header missing")
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	orgUserID, err := uuid.Parse(orgUserIDString)
	if err != nil {
		h.log.Error("Invalid X-Vetchi-OrgUserID", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	costCenterID, err := h.db.CreateCostCenter(r.Context(), db.CCenterReq{
		Name:      addCostCenterReq.Name,
		Notes:     addCostCenterReq.Notes,
		OrgUserID: orgUserID,
	})
	if err != nil {
		if errors.Is(err, db.ErrCostCenterAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	h.log.Debug("Created CostCenter", "ID", costCenterID, "orgUser", orgUserID)

	err = json.NewEncoder(w).Encode(vetchi.AddCostCenterResponse{
		Name: addCostCenterReq.Name,
	})
	if err != nil {
		h.log.Error("failed to encode response", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
