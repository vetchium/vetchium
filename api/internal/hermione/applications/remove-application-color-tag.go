package applications

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func RemoveApplicationColorTag(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered RemoveApplicationColorTag")
		var rmApplicationColorTagReq vetchi.RemoveApplicationColorTagRequest
		err := json.NewDecoder(r.Body).Decode(&rmApplicationColorTagReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		h.Dbg("RemoveApplicationColorTag", "req", rmApplicationColorTagReq)

		if !h.Vator().Struct(w, &rmApplicationColorTagReq) {
			h.Dbg("invalid request", "req", rmApplicationColorTagReq)
			return
		}

		err = h.DB().RemoveApplicationColorTag(
			r.Context(),
			rmApplicationColorTagReq,
		)
		if err != nil {
			if errors.Is(err, db.ErrNoApplication) {
				h.Dbg("not found", "id", rmApplicationColorTagReq.ApplicationID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrApplicationStateInCompatible) {
				h.Dbg("state", "id", rmApplicationColorTagReq.ApplicationID)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to remove application color tag", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("removed application color tag")
		w.WriteHeader(http.StatusOK)
	}
}
