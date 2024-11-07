package costcenter

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func DefunctCostCenter(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DefunctCostCenter")
		var defunctCostCenterRequest vetchi.DefunctCostCenterRequest
		err := json.NewDecoder(r.Body).Decode(&defunctCostCenterRequest)
		if err != nil {
			h.Dbg("failed to decode defunct cost center request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &defunctCostCenterRequest) {
			h.Dbg("validation failed", "defunctCCReq", defunctCostCenterRequest)
			return
		}

		h.Dbg("validated", "defunctCostCenterRequest", defunctCostCenterRequest)

		err = h.DB().DefunctCostCenter(
			r.Context(),
			defunctCostCenterRequest,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("CC not found", "name", defunctCostCenterRequest.Name)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to defunct cost center", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("defuncted cost center", "defunctCCReq", defunctCostCenterRequest)
		w.WriteHeader(http.StatusOK)
	}
}
