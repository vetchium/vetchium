package applications

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func AddApplicationColorTag(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddApplicationColorTag")
		var addApplicationColorTagReq vetchi.AddApplicationColorTagRequest
		err := json.NewDecoder(r.Body).Decode(&addApplicationColorTagReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addApplicationColorTagReq) {
			h.Dbg("invalid request", "error", err)
			return
		}
		h.Dbg("validated", "req", addApplicationColorTagReq)

		err = h.DB().
			AddApplicationColorTag(r.Context(), addApplicationColorTagReq)
		if err != nil {
			if errors.Is(err, db.ErrApplicationStateInCompatible) {
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			if errors.Is(err, db.ErrNoApplication) {
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to add application color tag", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added application color tag")
		w.WriteHeader(http.StatusOK)
	}
}
