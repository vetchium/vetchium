package colleagues

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func FilterColleagues(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var filterColleaguesReq hub.FilterColleaguesRequest
		err := json.NewDecoder(r.Body).Decode(&filterColleaguesReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterColleaguesReq) {
			h.Dbg("Invalid request", "request", filterColleaguesReq)
			return
		}
		h.Dbg("validated", "filterColleaguesReq", filterColleaguesReq)

		colleagues, err := h.DB().FilterColleagues(
			r.Context(),
			filterColleaguesReq,
		)
		if err != nil {
			h.Dbg("failed to filter colleagues", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("filtered", "colleagues", colleagues)

		err = json.NewEncoder(w).Encode(colleagues)
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
