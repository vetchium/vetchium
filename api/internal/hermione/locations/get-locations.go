package locations

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/vhandler"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func GetLocations(h vhandler.VHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getLocationsReq vetchi.GetLocationsRequest
		err := json.NewDecoder(r.Body).Decode(&getLocationsReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getLocationsReq) {
			h.Log().Error("failed to validate getLocationsReq", "error", err)
			return
		}

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		states := []string{}
		for _, state := range getLocationsReq.States {
			// already validated by vator
			states = append(states, string(state))
		}

		if getLocationsReq.Limit == 0 {
			getLocationsReq.Limit = 100
		}

		locations, err := h.DB().GetLocations(r.Context(), db.GetLocationsReq{
			States:        states,
			EmployerID:    orgUser.EmployerID,
			PaginationKey: getLocationsReq.PaginationKey,
			Limit:         getLocationsReq.Limit,
		})
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(locations)
		if err != nil {
			h.Log().Error("failed to encode locations", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
