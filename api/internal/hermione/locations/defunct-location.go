package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func DefunctLocation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered DefunctLocation")
		var defunctLocationRequest employer.DefunctLocationRequest
		err := json.NewDecoder(r.Body).Decode(&defunctLocationRequest)
		if err != nil {
			h.Dbg("failed to decode defunct location request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &defunctLocationRequest) {
			h.Dbg("validation failed", "defunctLocReq", defunctLocationRequest)
			return
		}
		h.Dbg("validated", "defunctLocationReq", defunctLocationRequest)

		err = h.DB().DefunctLocation(r.Context(), defunctLocationRequest)
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				h.Dbg("not found", "title", defunctLocationRequest.Title)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to defunct location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("defuncted location", "defunctLocReq", defunctLocationRequest)
		w.WriteHeader(http.StatusOK)
	}
}
