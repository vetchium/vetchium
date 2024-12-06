package applications

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func MyApplications(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered MyApplications")
		var myApplicationsReq vetchi.MyApplicationsRequest
		err := json.NewDecoder(r.Body).Decode(&myApplicationsReq)
		if err != nil {
			h.Dbg("failed to decode my applications request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &myApplicationsReq) {
			h.Dbg("validation failed", "myApplicationsReq", myApplicationsReq)
			return
		}

		h.Dbg("validated", "myApplicationsReq", myApplicationsReq)

		hubApplications, err := h.DB().
			MyApplications(r.Context(), myApplicationsReq)
		if err != nil {
			h.Dbg("failed to get my applications", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
		h.Dbg("got my applications", "hubApplications", hubApplications)

		err = json.NewEncoder(w).Encode(hubApplications)
		if err != nil {
			h.Dbg("failed to encode my applications", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
