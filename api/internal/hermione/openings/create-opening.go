package openings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
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

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		openingID, err := h.DB().CreateOpening(r.Context(), db.CreateOpeningReq{
			Title:              createOpeningReq.Title,
			Positions:          createOpeningReq.Positions,
			JD:                 createOpeningReq.JD,
			Recruiters:         createOpeningReq.Recruiters,
			HiringManager:      string(createOpeningReq.HiringManager),
			CostCenterName:     string(createOpeningReq.CostCenterName),
			EmployerNotes:      createOpeningReq.EmployerNotes,
			LocationTitles:     createOpeningReq.LocationTitles,
			RemoteCountryCodes: createOpeningReq.RemoteCountryCodes,
			RemoteTimezones:    createOpeningReq.RemoteTimezones,
			OpeningType:        string(createOpeningReq.OpeningType),
			YoeMin:             createOpeningReq.YoeMin,
			YoeMax:             createOpeningReq.YoeMax,
			MinEducationLevel:  createOpeningReq.MinEducationLevel,
			Salary:             createOpeningReq.Salary,
			EmployerID:         orgUser.EmployerID,
			CreatedBy:          orgUser.ID,
		})
		if err != nil {
			h.Dbg("failed to create opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created opening", "openingID", openingID)
		w.WriteHeader(http.StatusOK)
	}
}
