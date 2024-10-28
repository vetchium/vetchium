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

func UpdateLocation(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateLocationReq vetchi.UpdateLocationRequest
		err := json.NewDecoder(r.Body).Decode(&updateLocationReq)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateLocationReq) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
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
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
