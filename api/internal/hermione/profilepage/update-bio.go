package profilepage

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/hub"
)

func UpdateBio(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var updateBioRequest hub.UpdateBioRequest
		if err := json.NewDecoder(r.Body).Decode(&updateBioRequest); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &updateBioRequest) {
			h.Dbg("validation failed", "updateBioRequest", updateBioRequest)
			return
		}
		h.Dbg("validated", "updateBioRequest", updateBioRequest)

		if updateBioRequest.FullName == nil &&
			updateBioRequest.ShortBio == nil &&
			updateBioRequest.LongBio == nil {
			h.Dbg("no valid fields to update")
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		err := h.DB().UpdateBio(r.Context(), updateBioRequest)
		if err != nil {
			h.Err("failed to update bio", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("bio updated", "updateBioRequest", updateBioRequest)
		w.WriteHeader(http.StatusOK)
	}
}
