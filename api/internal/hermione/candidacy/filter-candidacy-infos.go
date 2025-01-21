package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func FilterCandidacyInfos(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterCandidacyInfos")

		var filterCandidacyInfosReq employer.FilterCandidacyInfosRequest
		err := json.NewDecoder(r.Body).Decode(&filterCandidacyInfosReq)
		if err != nil {
			h.Dbg("Error decoding request body", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterCandidacyInfosReq) {
			h.Dbg("Error validating request body")
			return
		}

		h.Dbg("validated", "filterCandidacyInfosReq", filterCandidacyInfosReq)

		if filterCandidacyInfosReq.Limit == 0 {
			filterCandidacyInfosReq.Limit = 40
		}

		candidaciesInfo, err := h.DB().
			FilterEmployerCandidacyInfos(r.Context(), filterCandidacyInfosReq)
		if err != nil {
			h.Dbg("Error getting candidacies info", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got candidacies info", "candidaciesInfo", candidaciesInfo)

		err = json.NewEncoder(w).Encode(candidaciesInfo)
		if err != nil {
			h.Dbg("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
