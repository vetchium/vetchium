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
		h.Log().Debug("Entered GetLocations")
		var getLocationsReq vetchi.GetLocationsRequest
		err := json.NewDecoder(r.Body).Decode(&getLocationsReq)
		if err != nil {
			h.Log().Debug("failed to decode getLocationsReq", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getLocationsReq) {
			h.Log().Debug("failed to validate getLocationsReq", "error", err)
			return
		}
		h.Log().Debug("Validated", "getLocationsReq", getLocationsReq)

		orgUser, ok := r.Context().Value(middleware.OrgUserCtxKey).(db.OrgUser)
		if !ok {
			h.Log().Error("failed to get orgUser from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Log().Debug("Got orgUser", "orgUser", orgUser)

		states := []string{}
		for _, state := range getLocationsReq.States {
			// already validated by vator
			states = append(states, string(state))
		}
		if len(states) == 0 {
			states = []string{string(vetchi.ActiveLocation)}
		}
		h.Log().Debug("States OK", "states", states)

		if getLocationsReq.Limit == 0 {
			getLocationsReq.Limit = 100
		}
		h.Log().Debug("Limit OK", "limit", getLocationsReq.Limit)

		locations, err := h.DB().GetLocations(r.Context(), db.GetLocationsReq{
			States:        states,
			EmployerID:    orgUser.EmployerID,
			PaginationKey: getLocationsReq.PaginationKey,
			Limit:         getLocationsReq.Limit,
		})
		if err != nil {
			h.Log().Debug("failed to get locations", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Log().Debug("Got locations", "locations", locations)
		err = json.NewEncoder(w).Encode(locations)
		if err != nil {
			h.Log().Error("failed to encode locations", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
