package incognito

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/util"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/api/pkg/vetchi"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

func AddIncognitoPost(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddIncognitoPost")
		var addIncognitoPostReq hub.AddIncognitoPostRequest
		if err := json.NewDecoder(r.Body).Decode(&addIncognitoPostReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addIncognitoPostReq) {
			h.Dbg(
				"validation failed",
				"addIncognitoPostReq",
				addIncognitoPostReq,
			)
			http.Error(w, "validation failed", http.StatusBadRequest)
			return
		}

		h.Dbg("Validated", "addIncognitoPostReq", addIncognitoPostReq)

		incognitoPostID := util.RandomUniqueID(vetchi.PostIDLenBytes)

		err := h.DB().AddIncognitoPost(r.Context(), db.AddIncognitoPostRequest{
			Context:             r.Context(),
			IncognitoPostID:     incognitoPostID,
			AddIncognitoPostReq: addIncognitoPostReq,
		})
		if err != nil {
			if err == db.ErrInvalidTagIDs {
				h.Dbg("invalid tag IDs provided", "error", err)
				w.WriteHeader(http.StatusBadRequest)
				err = json.NewEncoder(w).Encode(common.ValidationErrors{
					Errors: []string{"tags"},
				})
				if err != nil {
					h.Err("failed to encode validation errors", "error", err)
				}
				return
			}
			h.Dbg("adding incognito post failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(hub.AddIncognitoPostResponse{
			IncognitoPostID: incognitoPostID,
		})
		if err != nil {
			h.Dbg("encoding failed", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
