package openings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func CreateOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered CreateOpening")
		var createOpeningReq vetchi.CreateOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&createOpeningReq)
		if err != nil {
			h.Dbg("failed to decode create opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &createOpeningReq) {
			h.Dbg("validation failed", "createOpeningReq", createOpeningReq)
			return
		}
		h.Dbg("validated", "createOpeningReq", createOpeningReq)

		openingID, err := h.DB().CreateOpening(r.Context(), createOpeningReq)
		if err != nil {
			h.Dbg("failed to create opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("created opening", "openingID", openingID)
		w.WriteHeader(http.StatusOK)
	}
}
