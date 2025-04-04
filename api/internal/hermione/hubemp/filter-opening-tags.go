package hubemp

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
)

func FilterOpeningTags(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered FilterOpeningTags")

		var filterOpeningTagsReq common.FilterOpeningTagsRequest
		err := json.NewDecoder(r.Body).Decode(&filterOpeningTagsReq)
		if err != nil {
			h.Dbg("failed to decode filter opening tags request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &filterOpeningTagsReq) {
			h.Dbg(
				"validation failed",
				"filterOpeningTagsReq",
				filterOpeningTagsReq,
			)
			return
		}

		tags, err := h.DB().FilterOpeningTags(r.Context(), filterOpeningTagsReq)
		if err != nil {
			h.Err("failed to filter opening tags", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("filtered opening tags", "tags", tags)
		err = json.NewEncoder(w).Encode(tags)
		if err != nil {
			h.Err("failed to encode filter opening tags response", "error", err)
			return
		}
	}
}
