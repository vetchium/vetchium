package applications

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func GetApplications(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetApplications")
		var getApplicationsRequest vetchi.GetApplicationsRequest
		err := json.NewDecoder(r.Body).Decode(&getApplicationsRequest)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h.Dbg("GetApplications request", "request", getApplicationsRequest)

		if !h.Vator().Struct(w, &getApplicationsRequest) {
			h.Dbg("failed to validate request")
			return
		}
		h.Dbg("validated", "getApplicationsReq", getApplicationsRequest)

		applications, err := h.DB().
			GetApplicationsForEmployer(r.Context(), getApplicationsRequest)
		if err != nil {
			h.Dbg("failed to get applications", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("got applications", "applications", applications)

		err = json.NewEncoder(w).Encode(applications)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
