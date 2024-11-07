package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func UpdateOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateOpening")
		var updateOpeningReq vetchi.UpdateOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&updateOpeningReq)
		if err != nil {
			h.Dbg("failed to decode update opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateOpeningReq) {
			h.Dbg("validation failed", "updateOpeningReq", updateOpeningReq)
			return
		}
		h.Dbg("validated", "updateOpeningReq", updateOpeningReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().UpdateOpening(r.Context(), db.UpdateOpeningReq{
			ID:                 updateOpeningReq.ID,
			Title:              updateOpeningReq.Title,
			Positions:          updateOpeningReq.Positions,
			JD:                 updateOpeningReq.JD,
			Recruiters:         updateOpeningReq.Recruiters,
			HiringManager:      string(updateOpeningReq.HiringManager),
			CostCenterName:     string(updateOpeningReq.CostCenterName),
			EmployerNotes:      updateOpeningReq.EmployerNotes,
			LocationTitles:     updateOpeningReq.LocationTitles,
			RemoteCountryCodes: updateOpeningReq.RemoteCountryCodes,
			RemoteTimezones:    updateOpeningReq.RemoteTimezones,
			OpeningType:        string(updateOpeningReq.OpeningType),
			YoeMin:             updateOpeningReq.YoeMin,
			YoeMax:             updateOpeningReq.YoeMax,
			MinEducationLevel:  updateOpeningReq.MinEducationLevel,
			Salary:             updateOpeningReq.Salary,
			EmployerID:         orgUser.EmployerID,
			UpdatedBy:          orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", updateOpeningReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("updated opening", "id", updateOpeningReq.ID)
		w.WriteHeader(http.StatusOK)
	}
}
