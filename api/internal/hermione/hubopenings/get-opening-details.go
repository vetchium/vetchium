package hubopenings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func GetOpeningDetails(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetOpeningDetails")
		var getOpeningDetailsReq hub.GetHubOpeningDetailsRequest
		err := json.NewDecoder(r.Body).Decode(&getOpeningDetailsReq)
		if err != nil {
			h.Dbg("failed to decode get opening details request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getOpeningDetailsReq) {
			h.Dbg("validation failed", "req", getOpeningDetailsReq)
			return
		}
		h.Dbg("validated", "req", getOpeningDetailsReq)

		openingDetails, err := h.DB().
			GetHubOpeningDetails(r.Context(), getOpeningDetailsReq)
		if err != nil {
			h.Dbg("failed to get opening details", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got opening details", "openingDetails", openingDetails)
		err = json.NewEncoder(w).Encode(openingDetails)
		if err != nil {
			h.Err("failed to encode opening details", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
