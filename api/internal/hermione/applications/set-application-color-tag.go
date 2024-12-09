package applications

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func SetApplicationColorTag(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered SetApplicationColorTag")
		var setApplicationColorTagReq employer.SetApplicationColorTagRequest
		err := json.NewDecoder(r.Body).Decode(&setApplicationColorTagReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &setApplicationColorTagReq) {
			h.Dbg("invalid request", "error", err)
			return
		}
		h.Dbg("validated", "req", setApplicationColorTagReq)

		err = h.DB().
			SetApplicationColorTag(r.Context(), setApplicationColorTagReq)
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
