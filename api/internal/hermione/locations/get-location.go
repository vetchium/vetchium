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

func GetLocation(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getLocationReq vetchi.GetLocationRequest
		err := json.NewDecoder(r.Body).Decode(&getLocationReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getLocationReq) {
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Err("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		location, err := h.DB().GetLocByName(r.Context(), db.GetLocByNameReq{
			Title:      getLocationReq.Title,
			EmployerID: orgUser.EmployerID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(location)
		if err != nil {
			h.Err("failed to encode location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
