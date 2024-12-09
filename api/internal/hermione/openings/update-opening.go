package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func UpdateOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered UpdateOpening")
		var updateOpeningReq employer.UpdateOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&updateOpeningReq)
		if err != nil {
			h.Dbg("failed to decode update opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateOpeningReq) {
			h.Dbg("validation failed", "updateOpeningReq", updateOpeningReq)
			return
		}
		h.Dbg("validated", "updateOpeningReq", updateOpeningReq)

		err = h.DB().UpdateOpening(r.Context(), updateOpeningReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", updateOpeningReq.OpeningID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to update opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("updated opening", "id", updateOpeningReq.OpeningID)
		w.WriteHeader(http.StatusOK)
	}
}
