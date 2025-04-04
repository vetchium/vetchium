package colleagues

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func EndorseApplication(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EndorseApplication")
		var endorseReq hub.EndorseApplicationRequest
		if err := json.NewDecoder(r.Body).Decode(&endorseReq); err != nil {
			h.Dbg("Error decoding request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &endorseReq) {
			h.Dbg("Invalid request", "endorseReq", endorseReq)
			return
		}
		h.Dbg("Validated", "endorseReq", endorseReq)

		err := h.DB().EndorseApplication(r.Context(), endorseReq)
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("Application not found or not allowed", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}
			h.Dbg("Error endorsing application", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Endorsed application", "endorseReq", endorseReq)
	}
}
