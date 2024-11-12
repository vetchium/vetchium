package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func CreateOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered CreateOpening")
		var createOpeningReq vetchi.CreateOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&createOpeningReq)
		if err != nil {
			h.Dbg("failed to decode create opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &createOpeningReq) {
			h.Dbg("validation failed", "createOpeningReq", createOpeningReq)
			return
		}
		h.Dbg("validated", "createOpeningReq", createOpeningReq)

		if createOpeningReq.YoeMax < createOpeningReq.YoeMin {
			h.Dbg("yoe_max < yoe min", "createOpeningReq", createOpeningReq)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(vetchi.ValidationErrors{
				Errors: []string{"yoe_min", "yoe_max"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		if createOpeningReq.Salary != nil &&
			(createOpeningReq.Salary.MinAmount > createOpeningReq.Salary.MaxAmount) {
			h.Dbg("salary min > max", "createOpeningReq", createOpeningReq)
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(vetchi.ValidationErrors{
				Errors: []string{"salary"},
			})
			if err != nil {
				h.Err("failed to encode validation errors", "error", err)
			}
			return
		}

		openingID, err := h.DB().CreateOpening(r.Context(), createOpeningReq)
		if err != nil {
			if errors.Is(err, db.ErrNoRecruiter) ||
				errors.Is(err, db.ErrNoLocation) ||
				errors.Is(err, db.ErrNoHiringManager) ||
				errors.Is(err, db.ErrNoCostCenter) {
				h.Dbg("location or team or recruiter not found", "error", err)
				http.Error(w, err.Error(), http.StatusUnprocessableEntity)
				return
			}

			h.Err("failed to create opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created opening", "openingID", openingID)
		err = json.NewEncoder(w).Encode(vetchi.CreateOpeningResponse{
			OpeningID: openingID,
		})
		if err != nil {
			h.Err("failed to encode create opening response", "error", err)
			return
		}
	}
}
