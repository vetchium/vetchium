package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddLocation(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddLocation")
		var addLocationReq vetchi.AddLocationRequest
		err := json.NewDecoder(r.Body).Decode(&addLocationReq)
		if err != nil {
			h.Dbg("failed to decode add location request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addLocationReq) {
			h.Dbg("validation failed", "addLocationReq", addLocationReq)
			return
		}
		h.Dbg("validated", "addLocationReq", addLocationReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		locationID, err := h.DB().AddLocation(r.Context(), db.AddLocationReq{
			Title:            addLocationReq.Title,
			CountryCode:      addLocationReq.CountryCode,
			PostalAddress:    addLocationReq.PostalAddress,
			PostalCode:       addLocationReq.PostalCode,
			OpenStreetMapURL: addLocationReq.OpenStreetMapURL,
			CityAka:          addLocationReq.CityAka,
			EmployerID:       orgUser.EmployerID,
			OrgUserID:        orgUser.ID,
		})
		if err != nil {
			if errors.Is(err, db.ErrDupLocationName) {
				h.Dbg("location exists", "addLocationReq", addLocationReq)
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}

			h.Dbg("failed to add location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added location", "locationID", locationID)
		w.WriteHeader(http.StatusOK)
	}
}
