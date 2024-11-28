package hubopenings

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func ApplyForOpeningHandler(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered ApplyForOpeningHandler")
		var applyForOpeningReq vetchi.ApplyForOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&applyForOpeningReq)
		if err != nil {
			h.Dbg("failed to decode apply for opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &applyForOpeningReq) {
			h.Dbg("validation failed", "applyForOpeningReq", applyForOpeningReq)
			return
		}
		h.Dbg("validated", "applyForOpeningReq", applyForOpeningReq)

	}
}
