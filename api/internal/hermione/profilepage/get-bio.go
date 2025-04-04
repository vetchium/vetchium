package profilepage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func GetBio(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var getBioRequest hub.GetBioRequest
		if err := json.NewDecoder(r.Body).Decode(&getBioRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getBioRequest) {
			h.Dbg("validation failed", "getBioRequest", getBioRequest)
			return
		}
		h.Dbg("validated request", "getBioRequest", getBioRequest)

		bio, err := h.DB().GetBio(r.Context(), getBioRequest.Handle)
		if err != nil {
			if errors.Is(err, db.ErrNoHubUser) {
				h.Dbg("no hub user found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Err("failed to get bio", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Got bio", "bio", bio)
		err = json.NewEncoder(w).Encode(bio)
		if err != nil {
			h.Err("failed to encode bio", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
