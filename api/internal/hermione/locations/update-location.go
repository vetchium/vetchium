package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func UpdateLocation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateLocation")
		var updateLocationReq vetchi.UpdateLocationRequest
		err := json.NewDecoder(r.Body).Decode(&updateLocationReq)
		if err != nil {
			h.Dbg("failed to decode update location request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateLocationReq) {
			h.Dbg("validation failed", "updateLocationReq", updateLocationReq)
			return
		}
		h.Dbg("validated", "updateLocationReq", updateLocationReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUserTO)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		updateLocReq := db.UpdateLocationReq{
			Title:            updateLocationReq.Title,
			CountryCode:      updateLocationReq.CountryCode,
			PostalAddress:    updateLocationReq.PostalAddress,
			PostalCode:       updateLocationReq.PostalCode,
			OpenStreetMapURL: updateLocationReq.OpenStreetMapURL,
			CityAka:          updateLocationReq.CityAka,
			EmployerID:       orgUser.EmployerID,
		}

		err = h.DB().UpdateLocation(r.Context(), updateLocReq)
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				h.Dbg("not found", "title", updateLocationReq.Title)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("updated location", "updateLocReq", updateLocReq)
		w.WriteHeader(http.StatusOK)
	}
}
